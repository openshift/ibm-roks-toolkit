package kubeadminpwd

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/openshift/hypershift-toolkit/pkg/cmd/cpoperator"
	"github.com/openshift/hypershift-toolkit/pkg/controllers"
)

const (
	ManifestBootstrapperPod = "manifests-bootstrapper"
)

func Setup(cfg *cpoperator.ControlPlaneOperatorConfig) error {
	reconciler := &OAuthRestarter{
		Client:    cfg.Manager().GetClient(),
		Namespace: cfg.Namespace(),
		Log:       cfg.Logger().WithName("OAuthRestarter"),
	}
	c, err := controller.New("oauth-restarter", cfg.Manager(), controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}
	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}}, controllers.NamedResourceHandler(ManifestBootstrapperPod)); err != nil {
		return err
	}
	return nil
}
