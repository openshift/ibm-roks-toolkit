package openshiftapiserver

import (
	"bytes"
	"context"
	"crypto/md5" //#nosec - not used for security purposes
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kcache "k8s.io/apimachinery/pkg/util/cache"
	kubeinformers "k8s.io/client-go/informers"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/yaml"

	operatorv1 "github.com/openshift/api/operator/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/images"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/project"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

	"github.com/openshift/ibm-roks-toolkit/pkg/cmd/cpoperator"
	"github.com/openshift/ibm-roks-toolkit/pkg/controllers"
)

func Setup(cfg *cpoperator.ControlPlaneOperatorConfig) error {
	targetCfg := cfg.TargetConfig()
	kubeInformers := kubeinformers.NewSharedInformerFactoryWithOptions(cfg.TargetKubeClient(), controllers.DefaultResync, kubeinformers.WithNamespace("openshift-apiserver"))
	configClient, err := configclient.NewForConfig(targetCfg)
	if err != nil {
		return err
	}
	configInformers := configinformers.NewSharedInformerFactory(configClient, controllers.DefaultResync)
	operatorClient := newAPIServerOperatorClient(
		cfg.KubeClient(),
		cfg.Namespace(),
		cfg.Logger().WithName("OpenShiftAPIServerClient"),
	)

	recorder := events.NewLoggingEventRecorder("openshift-apiserver-observers")
	c := configobserver.NewConfigObserver(
		operatorClient,
		recorder,
		configobservation.Listers{
			ResourceSync:        &noopResourceSyncer{},
			APIServerLister_:    configInformers.Config().V1().APIServers().Lister(),
			ImageConfigLister:   configInformers.Config().V1().Images().Lister(),
			ProjectConfigLister: configInformers.Config().V1().Projects().Lister(),
			ProxyLister_:        configInformers.Config().V1().Proxies().Lister(),
			IngressConfigLister: configInformers.Config().V1().Ingresses().Lister(),
			SecretLister_:       kubeInformers.Core().V1().Secrets().Lister(),
			PreRunCachesSynced: []cache.InformerSynced{
				configInformers.Config().V1().Images().Informer().HasSynced,
				configInformers.Config().V1().Projects().Informer().HasSynced,
				configInformers.Config().V1().Proxies().Informer().HasSynced,
				configInformers.Config().V1().Ingresses().Informer().HasSynced,
			},
		},
		[]factory.Informer{},
		images.ObserveInternalRegistryHostname,
		images.ObserveAllowedRegistriesForImport,
		project.ObserveProjectRequestMessage,
		project.ObserveProjectRequestTemplateName,
	)
	cfg.Manager().Add(manager.RunnableFunc(func(ctx context.Context) error {
		configInformers.Start(ctx.Done())
		return nil
	}))
	cfg.Manager().Add(manager.RunnableFunc(func(ctx context.Context) error {
		c.Run(ctx, 1)
		return nil
	}))
	return nil
}

const (
	apiserverConfigMapName  = "openshift-apiserver-config"
	openshiftDeploymentName = "openshift-apiserver"
)

type apiServerOperatorClient struct {
	Client    kubeclient.Interface
	Namespace string
	Logger    logr.Logger

	configCache *kcache.Expiring
}

func newAPIServerOperatorClient(c kubeclient.Interface, ns string, logger logr.Logger) v1helpers.OperatorClient {
	return &apiServerOperatorClient{
		Client:      c,
		Namespace:   ns,
		Logger:      logger,
		configCache: kcache.NewExpiring(),
	}
}

func (c *apiServerOperatorClient) Informer() cache.SharedIndexInformer {
	panic("informer not supported")
}
func (c *apiServerOperatorClient) GetObjectMeta() (meta *metav1.ObjectMeta, err error) {
	panic("operator object meta not found")
}

var defaultExpirationTime = 24 * time.Hour

func (c *apiServerOperatorClient) GetOperatorState() (spec *operatorv1.OperatorSpec, status *operatorv1.OperatorStatus, resourceVersion string, err error) {
	var cm *corev1.ConfigMap

	cmObj, ok := c.configCache.Get("config")
	if !ok || cmObj == nil {
		cm, err = c.Client.CoreV1().ConfigMaps(c.Namespace).Get(context.TODO(), apiserverConfigMapName, metav1.GetOptions{})
		if err != nil {
			return
		}
		c.configCache.Set("config", cm, defaultExpirationTime)
	} else {
		cm, ok = cmObj.(*corev1.ConfigMap)
		if !ok {
			c.configCache.Delete("config")
			err = fmt.Errorf("unexpected object of type %T in cache", cmObj)
			return
		}
	}
	configYAML := []byte(cm.Data["config.yaml"])
	var configJSON []byte
	configJSON, err = yaml.YAMLToJSON(configYAML)
	if err != nil {
		return
	}
	configJSON, err = filterManagedConfigKeys(configJSON)
	if err != nil {
		return
	}
	spec = &operatorv1.OperatorSpec{}
	status = &operatorv1.OperatorStatus{}
	spec.ObservedConfig.Raw = configJSON
	resourceVersion = cm.ResourceVersion
	return
}

// UpdateOperatorSpec updates the spec of the operator, assuming the given resource version.
func (c *apiServerOperatorClient) UpdateOperatorSpec(ctx context.Context, oldResourceVersion string, in *operatorv1.OperatorSpec) (out *operatorv1.OperatorSpec, newResourceVersion string, err error) {
	var cm *corev1.ConfigMap
	cm, err = c.Client.CoreV1().ConfigMaps(c.Namespace).Get(ctx, apiserverConfigMapName, metav1.GetOptions{})
	if err != nil {
		return
	}
	if cm.ResourceVersion != oldResourceVersion {
		err = fmt.Errorf("resource version does not match")
		return
	}
	var updateJSON []byte
	updateJSON, err = in.ObservedConfig.MarshalJSON()
	if err != nil {
		return
	}
	var configBytes []byte
	configBytes, err = mergeConfig([]byte(cm.Data["config.yaml"]), updateJSON)
	if err != nil {
		return
	}
	cm.Data["config.yaml"] = string(configBytes)
	c.Logger.Info("Updating OpenShift APIServer configmap")
	c.configCache.Delete("config")
	_, err = c.Client.CoreV1().ConfigMaps(c.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		return
	}
	dataHash := calculateHash(configBytes)
	var deployment *appsv1.Deployment
	deployment, err = c.Client.AppsV1().Deployments(c.Namespace).Get(ctx, openshiftDeploymentName, metav1.GetOptions{})
	if err != nil {
		return
	}
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	c.Logger.Info("Updating OpenShift APIServer deployment")
	deployment.Spec.Template.ObjectMeta.Annotations["config-checksum"] = dataHash
	_, err = c.Client.AppsV1().Deployments(c.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	return
}

func mergeConfig(existingYAML, updateJSON []byte) (updatedYAML []byte, err error) {
	var existingJSON []byte
	existingJSON, err = yaml.YAMLToJSON(existingYAML)
	if err != nil {
		return
	}
	existingConfig := map[string]interface{}{}
	if err = json.NewDecoder(bytes.NewBuffer(existingJSON)).Decode(&existingConfig); err != nil {
		return
	}
	updateConfig := map[string]interface{}{}
	if err = json.NewDecoder(bytes.NewBuffer(updateJSON)).Decode(&updateConfig); err != nil {
		return
	}

	if value, exists := updateConfig["projectConfig"]; exists {
		existingConfig["projectConfig"] = value
	} else {
		delete(existingConfig, "projectConfig")
	}

	if value, exists := updateConfig["imagePolicyConfig"]; exists {
		var resultValue map[string]interface{}
		if existingValue, ok := existingConfig["imagePolicyConfig"]; !ok {
			resultValue = map[string]interface{}{}
		} else {
			resultValue = existingValue.(map[string]interface{})
		}
		if mapValue, ok := value.(map[string]interface{}); ok {
			for _, key := range []string{"internalRegistryHostname", "allowedRegistriesForImport"} {
				if value, exists := mapValue[key]; exists {
					resultValue[key] = value
				} else {
					delete(resultValue, key)
				}
			}
		}
		existingConfig["imagePolicyConfig"] = resultValue
	} else {
		delete(existingConfig, "imagePolicyConfig")
	}

	var mergedConfig []byte
	mergedConfig, err = json.Marshal(existingConfig)
	if err != nil {
		return
	}

	updatedYAML, err = yaml.JSONToYAML(mergedConfig)
	return
}

// filterManagedConfigKeys returns JSON that contains only the keys managed by the
// observer controller from a bigger config JSON
func filterManagedConfigKeys(in []byte) (out []byte, err error) {
	inputConfig := map[string]interface{}{}
	if err = json.NewDecoder(bytes.NewBuffer(in)).Decode(&inputConfig); err != nil {
		return
	}
	outputConfig := map[string]interface{}{}
	for key := range inputConfig {
		switch key {
		case "projectConfig":
			outputConfig[key] = inputConfig[key]
		case "imagePolicyConfig":
			resultValue := map[string]interface{}{}
			if mapValue, ok := inputConfig[key].(map[string]interface{}); ok {
				for key2 := range mapValue {
					switch key2 {
					case "internalRegistryHostname", "allowedRegistriesForImport":
						resultValue[key2] = mapValue[key2]
					}
				}
			}
			outputConfig[key] = resultValue
		}
	}
	out, err = json.Marshal(outputConfig)
	return
}

// UpdateOperatorStatus updates the status of the operator, assuming the given resource version.
func (c *apiServerOperatorClient) UpdateOperatorStatus(ctx context.Context, oldResourceVersion string, in *operatorv1.OperatorStatus) (out *operatorv1.OperatorStatus, err error) {
	return
}

type noopResourceSyncer struct {
}

func (*noopResourceSyncer) SyncConfigMap(destination, source resourcesynccontroller.ResourceLocation) error {
	panic("configmap sync requested")
}

func (*noopResourceSyncer) SyncSecret(destination, source resourcesynccontroller.ResourceLocation) error {
	panic("secret sync requested")
}

func calculateHash(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b)) //#nosec - not used for security purposes
}
