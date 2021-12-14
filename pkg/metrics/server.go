package metrics

import (
	"net/http"
	"time"

	"github.com/blang/semver"
	"github.com/prometheus/client_golang/prometheus"

	"k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/tools/clientcmd"
	basemetrics "k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/klog/v2"

	buildv1client "github.com/openshift/client-go/build/clientset/versioned"
	buildv1informers "github.com/openshift/client-go/build/informers/externalversions"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configv1informers "github.com/openshift/client-go/config/informers/externalversions"
	"github.com/openshift/cluster-kube-apiserver-operator/pkg/operator/configmetrics"
	buildmetrics "github.com/openshift/openshift-controller-manager/pkg/build/metrics/prometheus"
)

type Server struct {
	ListenAddress string
	CertFile      string
	KeyFile       string
}

func (s *Server) Run() error {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	cfg, err := kubeconfig.ClientConfig()
	if err != nil {
		klog.Errorf("Cannot get cluster configuration: %v", err)
		return err
	}

	configClient, err := configv1client.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("Cannot get config client: %v", err)
		return err
	}

	buildClient, err := buildv1client.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("Cannot get build client: %v", err)
		return err
	}

	configInformers := configv1informers.NewSharedInformerFactory(configClient, 30*time.Minute)
	buildInformers := buildv1informers.NewSharedInformerFactory(buildClient, 30*time.Minute)

	configmetrics.Register(configInformers)
	// Build metrics are registered to the prometheus default registry. This changes the
	// default registerer to the Kube legacy registry so that both configmetrics and build
	// metrics are registered to the same registry.
	prometheus.DefaultRegisterer = &legacyRegistryAdapter{}
	buildmetrics.IntializeMetricsCollector(buildInformers.Build().V1().Builds().Lister())

	done := server.SetupSignalHandler()
	configInformers.Start(done)
	buildInformers.Start(done)

	smux := http.NewServeMux()
	smux.Handle("/metrics", legacyregistry.HandlerWithReset())
	server := &http.Server{
		Addr:    s.ListenAddress,
		Handler: smux,
	}
	go func() {
		klog.Infof("Listening on %s", s.ListenAddress)
		<-done
		klog.Info("Done received, closing server")
		server.Close()
	}()
	return server.ListenAndServeTLS(s.CertFile, s.KeyFile)
}

type legacyRegistryAdapter struct{}

func (a *legacyRegistryAdapter) Register(c prometheus.Collector) error {
	return legacyregistry.Register(&registerableAdapter{Collector: c})
}

func (a *legacyRegistryAdapter) MustRegister(collectors ...prometheus.Collector) {
	registerable := make([]basemetrics.Registerable, len(collectors))
	for i, c := range collectors {
		registerable[i] = &registerableAdapter{Collector: c}
	}
	legacyregistry.MustRegister(registerable...)
}

func (a *legacyRegistryAdapter) Unregister(c prometheus.Collector) bool {
	return legacyregistry.DefaultGatherer.(basemetrics.KubeRegistry).Unregister(c)
}

type registerableAdapter struct {
	prometheus.Collector
}

func (a *registerableAdapter) Create(version *semver.Version) bool {
	return true
}

func (a *registerableAdapter) ClearState() {
}

func (a *registerableAdapter) FQName() string {
	return "ibm_roks_metrics_controller"
}
