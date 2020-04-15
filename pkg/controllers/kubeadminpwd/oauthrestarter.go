package kubeadminpwd

import (
	"context"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	OAuthDeploymentName = "oauth-openshift"
)

type OAuthRestarter struct {
	// Client is a client of the management cluster
	client.Client

	// Log is the logger for this controller
	Log logr.Logger

	// Namespace is the namespace where the control plane of the cluster
	// lives on the management server
	Namespace string
}

func (o *OAuthRestarter) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	controllerLog := o.Log.WithValues("pod", req.NamespacedName.String())
	ctx := context.Background()

	// Ignore any namespace that is not the Namespace.
	if req.Namespace != o.Namespace {
		return ctrl.Result{}, nil
	}

	// Ignore all pods except the manifest bootstrapper
	if req.Name != ManifestBootstrapperPod {
		return ctrl.Result{}, nil
	}

	controllerLog.Info("Begin reconciling")

	pod := &corev1.Pod{}
	if err := o.Get(ctx, req.NamespacedName, pod); err != nil {
		return ctrl.Result{}, err
	}

	if !isComplete(pod) {
		controllerLog.Info("Pod has not yet completed")
		return ctrl.Result{}, nil
	}

	controllerLog.Info("Pod has completed, annotating the OAuth deployment")

	oauthDeployment := &appsv1.Deployment{}
	if err := o.Get(ctx, types.NamespacedName{Namespace: o.Namespace, Name: OAuthDeploymentName}, oauthDeployment); err != nil {
		return ctrl.Result{}, err
	}
	if oauthDeployment.Spec.Template.ObjectMeta.Annotations == nil {
		oauthDeployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
	}
	oauthDeployment.Spec.Template.ObjectMeta.Annotations["bootstrap-pod-resource-version"] = pod.ResourceVersion

	if err := o.Update(ctx, oauthDeployment); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func isComplete(pod *corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodSucceeded

}
