package clusteroperator

import (
	"context"
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/client-go/config/clientset/versioned/fake"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	common "github.com/openshift/ibm-roks-toolkit/pkg/controllers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestReconcile(t *testing.T) {
	versions := map[string]string{
		"release":    "alpha",
		"kubernetes": "beta",
	}

	testCases := []struct {
		testCase    string
		existingCOs []runtime.Object
	}{
		{
			testCase: "when there's no cluster operators",
		},
		{
			testCase: "when a clusterOperator exists with no status",
			existingCOs: []runtime.Object{
				&configv1.ClusterOperator{
					ObjectMeta: metav1.ObjectMeta{
						Name: "openshift-apiserver",
					},
				},
			},
		},
		{
			testCase: "when a clusterOperator exists with outdated versions",
			existingCOs: []runtime.Object{
				&configv1.ClusterOperator{
					ObjectMeta: metav1.ObjectMeta{
						Name: "openshift-apiserver",
					},
					Status: configv1.ClusterOperatorStatus{
						Conditions: nil,
						Versions: []configv1.OperandVersion{
							{
								Name:    "openshift-apiserver",
								Version: "outdated",
							},
							{
								Name:    "operator",
								Version: "outdated",
							},
						},
						RelatedObjects: nil,
						Extension:      runtime.RawExtension{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset(tc.existingCOs...)
			informerFactory := configinformers.NewSharedInformerFactory(fakeClient, common.DefaultResync)

			stopCh := make(chan struct{})
			informer := informerFactory.Config().V1().ClusterOperators().Informer()
			go informerFactory.Start(stopCh)
			cache.WaitForCacheSync(stopCh, informer.HasSynced)

			for _, obj := range tc.existingCOs {
				informerFactory.Config().V1().ClusterOperators().Informer().GetStore().Add(obj)
			}

			r := &ControlPlaneClusterOperatorSyncer{
				Versions: versions,
				Client:   fakeClient,
				Lister:   informerFactory.Config().V1().ClusterOperators().Lister(),
				Log:      ctrl.Log.WithName("TestReconcile"),
			}

			request := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name: "openshift-apiserver",
				},
			}

			var reconcileErr error
			err := wait.Poll(1*time.Second, 5*time.Second, func() (bool, error) {
				_, err := r.Reconcile(context.Background(), request)
				if err != nil {
					reconcileErr = err
					return false, nil
				}
				return true, nil
			})
			if err != nil {
				t.Errorf("failed to reconcile: %v: %v", err, reconcileErr)
			}

			gotClusterOperators, err := r.Client.ConfigV1().ClusterOperators().List(context.Background(), metav1.ListOptions{})
			if err != nil {
				t.Fatal(err)
			}

			if err := validateClusterOperators(gotClusterOperators.Items, versions); err != nil {
				t.Error(err)
			}

			stopCh <- struct{}{}
		})
	}
}

func validateClusterOperators(gotClusterOperators []configv1.ClusterOperator, versions map[string]string) error {
	notFound := sets.NewString(clusterOperatorNames.List()...)
	for _, co := range gotClusterOperators {
		if notFound.Has(co.Name) {
			notFound.Delete(co.Name)

			// find coInfo for current got cluster operator
			foundCOInfo := false
			for _, coInfo := range clusterOperators {
				if coInfo.Name == co.Name {
					foundCOInfo = true

					// find each expected operand VersionMapping in co.Status.Versions
					for operand, target := range coInfo.VersionMapping {
						found := false
						for _, operandVersion := range co.Status.Versions {
							if operandVersion.Name == operand {
								found = true
								if operandVersion.Version != versions[target] {
									return fmt.Errorf("operandVersion for %v does not match version %v", operand, operandVersion.Version)
								}
							}
						}
						if !found {
							return fmt.Errorf("operand %v not found in %v", operand, co.Status.Versions)
						}
					}
					break
				}
			}
			if !foundCOInfo {
				return fmt.Errorf("not coInfo found for %v", co.Name)
			}
		}
	}
	if notFound.Len() != 0 {
		return fmt.Errorf("some operators does not exist: %v", notFound)
	}
	return nil
}
