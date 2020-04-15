package clusterversion

import (
	"fmt"

	"github.com/go-logr/logr"

	ctrl "sigs.k8s.io/controller-runtime"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	configlister "github.com/openshift/client-go/config/listers/config/v1"
)

const (
	// NOTE: This needs to be set appropriately for the release associated with this code
	DefaultChannel = "stable-4.3"
)

type ClusterVersionReconciler struct {
	Client configclient.Interface
	Lister configlister.ClusterVersionLister
	Log    logr.Logger
}

func (r *ClusterVersionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	clusterVersion, err := r.Lister.Get(req.Name)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot fetch cluster version %s: %v", req.Name, err)
	}
	updateNeeded := false
	// Always default to default channel
	if clusterVersion.Spec.Channel != DefaultChannel {
		clusterVersion.Spec.Channel = DefaultChannel
		updateNeeded = true
	}
	// Remove any attempt at changing the clusterVersion
	if clusterVersion.Spec.DesiredUpdate != nil {
		clusterVersion.Spec.DesiredUpdate = nil
		updateNeeded = true
	}
	if updateNeeded {
		r.Log.Info("Updating clusterversion resource to desired values")
		_, err := r.Client.ConfigV1().ClusterVersions().Update(clusterVersion)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
