package openshiftapiserver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ctrl "sigs.k8s.io/controller-runtime"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	fakeconfigclient "github.com/openshift/client-go/config/clientset/versioned/fake"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/images"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation/project"
	"github.com/openshift/ibm-roks-toolkit/pkg/controllers/testconfigobserver"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
)

func TestControllerSync(t *testing.T) {
	tests := []struct {
		initial  string
		expected string
		config   []runtime.Object
	}{
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
projectConfig:
  projectRequestMessage: this is a request message
  projectRequestTemplate: openshift-config/the-template
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{
				&configv1.Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProjectSpec{
						ProjectRequestMessage: "this is a request message",
						ProjectRequestTemplate: configv1.TemplateReference{
							Name: "the-template",
						},
					},
				},
			},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
projectConfig:
  projectRequestMessage: this is a request message
  projectRequestTemplate: openshift-config/the-template
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
projectConfig:
  projectRequestMessage: ""
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{
				&configv1.Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProjectSpec{},
				},
			},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
projectConfig:
  projectRequestMessage: ""
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{
				&configv1.Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProjectSpec{},
				},
			},
		},
		{
			initial: `
apiVersion: kubecontrolplane.config.openshift.io/v1
kind: KubeControllerManagerConfig
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: kubecontrolplane.config.openshift.io/v1
imagePolicyConfig:
  internalRegistryHostname: test.internal.registry.host
kind: KubeControllerManagerConfig
projectConfig:
  projectRequestMessage: ""
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{
				&configv1.Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProjectSpec{},
				},
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.ImageStatus{
						InternalRegistryHostname: "test.internal.registry.host",
					},
				},
			},
		},
		{
			initial: `
apiVersion: kubecontrolplane.config.openshift.io/v1
kind: KubeControllerManagerConfig
imagePolicyConfig:
  internalRegistryHostname: test.internal.registry.host
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: kubecontrolplane.config.openshift.io/v1
kind: KubeControllerManagerConfig
projectConfig:
  projectRequestMessage: ""
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{
				&configv1.Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProjectSpec{},
				},
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.ImageStatus{},
				},
			},
		},
		{
			initial: `
apiVersion: kubecontrolplane.config.openshift.io/v1
imagePolicyConfig:
  allowedRegistriesForImport:
  - domainName: foo
  - domainName: bar
  internalRegistryHostname: test.internal.registry.host
kind: KubeControllerManagerConfig
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
`,
			expected: `
apiVersion: kubecontrolplane.config.openshift.io/v1
imagePolicyConfig:
  allowedRegistriesForImport:
  - domainName: foo
  - domainName: bar
kind: KubeControllerManagerConfig
projectConfig:
  projectRequestMessage: ""
serviceServingCert:
  certFile: /etc/kubernetes/config/service-ca.crt
			`,
			config: []runtime.Object{
				&configv1.Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ProjectSpec{},
				},
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.ImageSpec{
						AllowedRegistriesForImport: []configv1.RegistryLocation{
							{
								DomainName: "foo",
								Insecure:   false,
							},
							{
								DomainName: "bar",
								Insecure:   false,
							},
						},
					},
					Status: configv1.ImageStatus{},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			testSync(t, test.initial, test.expected, test.config...)
		})
	}
}

const testNamespace = "test-namespace"

func testSync(t *testing.T, initialConfig string, expectedConfig string, guestClusterConfig ...runtime.Object) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiserverConfigMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"config.yaml": initialConfig,
		},
	}
	oapiDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openshift-apiserver",
			Namespace: testNamespace,
		},
	}

	fakeClient := fake.NewSimpleClientset(configMap, oapiDeployment)

	operatorClient := newAPIServerOperatorClient(
		fakeClient,
		"test-namespace",
		ctrl.Log.WithName("test-controller-sync"),
	)

	configClient := fakeconfigclient.NewSimpleClientset(guestClusterConfig...)
	recorder := events.NewLoggingEventRecorder("openshift-apiserver-observers")
	syncContext := factory.NewSyncContext("openshift-apiserver-observers", recorder)

	configObserver := testconfigobserver.NewConfigObserver(
		operatorClient,
		configobservation.Listers{
			ResourceSync:        &noopResourceSyncer{},
			APIServerLister_:    nil,
			ImageConfigLister:   &imageLister{client: configClient},
			ProjectConfigLister: &projectLister{client: configClient},
			ProxyLister_:        nil,
			IngressConfigLister: nil,
			SecretLister_:       nil,
			PreRunCachesSynced:  nil,
		},
		[]factory.Informer{},
		images.ObserveInternalRegistryHostname,
		images.ObserveAllowedRegistriesForImport,
		project.ObserveProjectRequestMessage,
		project.ObserveProjectRequestTemplateName,
	)
	err := configObserver.Sync(context.Background(), syncContext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	configMap, err = fakeClient.CoreV1().ConfigMaps(testNamespace).Get(context.Background(), apiserverConfigMapName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(configMap.Data["config.yaml"]) != strings.TrimSpace(expectedConfig) {
		t.Errorf("unexpected config:\n%s\nexpected:\n%s", configMap.Data["config.yaml"], expectedConfig)
	}
}

type imageLister struct {
	client configclient.Interface
}

func (l *imageLister) List(selector labels.Selector) (ret []*configv1.Image, err error) {
	return []*configv1.Image{}, nil
}
func (l *imageLister) Get(name string) (*configv1.Image, error) {
	return l.client.ConfigV1().Images().Get(context.Background(), name, metav1.GetOptions{})
}

type projectLister struct {
	client configclient.Interface
}

func (l *projectLister) List(selector labels.Selector) (ret []*configv1.Project, err error) {
	return []*configv1.Project{}, nil
}
func (l *projectLister) Get(name string) (*configv1.Project, error) {
	return l.client.ConfigV1().Projects().Get(context.Background(), name, metav1.GetOptions{})
}
