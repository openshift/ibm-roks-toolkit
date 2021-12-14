package kubeletservingca

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/openshift/ibm-roks-toolkit/pkg/cmd/cpoperator"
)

const (
	ManagedConfigNamespace = "openshift-config-managed"
)

func Setup(cfg *cpoperator.ControlPlaneOperatorConfig) error {

	informerFactory := cfg.TargetKubeInformersForNamespace(ManagedConfigNamespace)
	configMaps := informerFactory.Core().V1().ConfigMaps()

	reconciler := &Syncer{
		InitialCA:    cfg.InitialCA(),
		TargetClient: cfg.TargetKubeClient(),
		Log:          cfg.Logger().WithName("KubeletServingCA"),
	}
	c, err := controller.New("kubelet-serving-ca", cfg.Manager(), controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}
	if err := c.Watch(&source.Informer{Informer: configMaps.Informer()}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}
	return nil
}
