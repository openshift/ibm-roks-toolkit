package openshiftcontrollermanager

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
	"github.com/openshift/cluster-openshift-controller-manager-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-openshift-controller-manager-operator/pkg/operator/configobservation/builds"
	"github.com/openshift/cluster-openshift-controller-manager-operator/pkg/operator/configobservation/images"
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
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
			`,
			config: []runtime.Object{},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
dockerPullSecret:
  internalRegistryHostname: internal.registry.host
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			config: []runtime.Object{
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.ImageStatus{
						InternalRegistryHostname: "internal.registry.host",
					},
				},
			},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
dockerPullSecret:
  internalRegistryHostname: internal.registry.host
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			config: []runtime.Object{
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
apiVersion: openshiftcontrolplane.config.openshift.io/v1
dockerPullSecret:
  internalRegistryHostname: internal.registry.host
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
build:
  buildDefaults:
    resources: {}
dockerPullSecret:
  internalRegistryHostname: internal.registry.host
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			config: []runtime.Object{
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.ImageStatus{
						InternalRegistryHostname: "internal.registry.host",
					},
				},
				&configv1.Build{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
				},
			},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
dockerPullSecret:
  internalRegistryHostname: internal.registry.host
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
build:
  buildDefaults:
    env:
    - name: FOO
      value: BAR
    resources: {}
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			config: []runtime.Object{
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.ImageStatus{},
				},
				&configv1.Build{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.BuildSpec{
						BuildDefaults: configv1.BuildDefaults{
							Env: []corev1.EnvVar{
								{
									Name:  "FOO",
									Value: "BAR",
								},
							},
						},
					},
				},
			},
		},
		{
			initial: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
build:
  buildDefaults:
    env:
    - name: FOO
      value: BAR
    resources: {}
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			expected: `
apiVersion: openshiftcontrolplane.config.openshift.io/v1
build:
  buildDefaults:
    resources: {}
ingress:
  ingressIPNetworkCIDR: ""
kind: OpenShiftControllerManagerConfig
kubeClientConfig:
  kubeConfig: /etc/kubernetes/secret/kubeconfig
`,
			config: []runtime.Object{
				&configv1.Image{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Status: configv1.ImageStatus{},
				},
				&configv1.Build{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster",
					},
					Spec: configv1.BuildSpec{},
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
			Name:      configMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"config.yaml": initialConfig,
		},
	}
	oapiDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: testNamespace,
		},
	}

	fakeClient := fake.NewSimpleClientset(configMap, oapiDeployment)

	operatorClient := newCMOperatorClient(
		fakeClient,
		testNamespace,
		ctrl.Log.WithName("test-controller-sync"),
	)

	configClient := fakeconfigclient.NewSimpleClientset(guestClusterConfig...)
	recorder := events.NewLoggingEventRecorder("openshift-controller-manager-observers")
	syncContext := factory.NewSyncContext("openshift-controller-manager-observers", recorder)

	configObserver := testconfigobserver.NewConfigObserver(
		operatorClient,
		configobservation.Listers{
			ImageConfigLister: &imageLister{client: configClient},
			BuildConfigLister: &buildLister{client: configClient},
		},
		[]factory.Informer{},
		images.ObserveInternalRegistryHostname,
		builds.ObserveBuildControllerConfig,
	)
	err := configObserver.Sync(context.Background(), syncContext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	configMap, err = fakeClient.CoreV1().ConfigMaps(testNamespace).Get(context.Background(), configMapName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a, e := strings.TrimSpace(configMap.Data["config.yaml"]), strings.TrimSpace(expectedConfig); a != e {
		t.Errorf("unexpected config:\n%s\nexpected:\n%s", a, e)
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

type buildLister struct {
	client configclient.Interface
}

func (l *buildLister) List(selector labels.Selector) (ret []*configv1.Build, err error) {
	return []*configv1.Build{}, nil
}
func (l *buildLister) Get(name string) (*configv1.Build, error) {
	return l.client.ConfigV1().Builds().Get(context.Background(), name, metav1.GetOptions{})
}
