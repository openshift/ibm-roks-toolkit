package clusteroperator

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/ibm-roks-toolkit/pkg/cmd/cpoperator"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

var (
	cfg          *rest.Config
	k8sClient    kubeclient.Interface
	configClient configclient.Interface
	testEnv      *envtest.Environment
	doneMgr      = make(chan struct{})
	ctx          = context.Background()
)

func TestClusterOperator(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Control Plane Operator (Cluster Operator Controller) Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("../", "../", "../", "test-config", "clusteroperator.crd.yaml")},
	}

	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	By("Setting up a new ControlPlaneOperatorConfig with a cluster-operator controller")
	versions := map[string]string{
		"release":    "alpha",
		"kubernetes": "beta",
	}
	controllerFuncs := map[string]cpoperator.ControllerSetupFunc{
		"cluster-operator": Setup,
	}

	cpoConfig := cpoperator.NewControlPlaneOperatorConfigWithRestConfig(cfg,
		cfg, "", nil, versions, []string{"cluster-operator"}, controllerFuncs)

	k8sClient = cpoConfig.TargetKubeClient()
	configClient = cpoConfig.TargetConfigClient()

	By("Creating the namespace")
	namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: cpoConfig.TargetNamespace()}}
	_, err = k8sClient.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	Expect(err).ToNot(HaveOccurred())

	By("Starting the manager")
	go func() {
		defer GinkgoRecover()
		Expect(cpoConfig.Start(ctx)).To(Succeed())
	}()

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	close(doneMgr)
	Expect(testEnv.Stop()).To(Succeed())
})

var _ = Describe("When running a clusterOperator controller", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	It("Should create all clusterOperator resources", func() {
		co := &configv1.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: "openshift-apiserver",
			},
		}
		By("creating a clusterOperator resource explicitly")
		co, err := configClient.ConfigV1().ClusterOperators().Create(ctx, co, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		By("waiting for all other control plane clusterOperator resources to be reconciled")
		Eventually(func() (bool, error) {
			gotClusterOperators, err := configClient.ConfigV1().ClusterOperators().List(ctx, metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())

			versions := map[string]string{
				"release":    "alpha",
				"kubernetes": "beta",
			}
			if err := validateClusterOperators(gotClusterOperators.Items, versions); err != nil {
				return false, nil
			}
			return true, nil
		}, 5*time.Second).Should(BeTrue())
	})
})
