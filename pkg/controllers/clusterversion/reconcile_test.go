package clusterversion

import (
	"context"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/client-go/config/clientset/versioned/fake"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	common "github.com/openshift/ibm-roks-toolkit/pkg/controllers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestReconcile(t *testing.T) {
	testCases := []struct {
		testCase   string
		existingCV *configv1.ClusterVersion
	}{
		{
			testCase: "when clusterVersion has no values",
			existingCV: &configv1.ClusterVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name: "version",
				},
			},
		},
		{
			testCase: "when clusterVersion has unexpected values",
			existingCV: &configv1.ClusterVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name: "version",
				},
				Spec: configv1.ClusterVersionSpec{
					ClusterID: "",
					DesiredUpdate: &configv1.Update{
						Version: "anything",
						Image:   "anything",
						Force:   false,
					},
					Upstream: "anything",
					Channel:  "anything",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset(tc.existingCV)
			informerFactory := configinformers.NewSharedInformerFactory(fakeClient, common.DefaultResync)
			stopCh := make(chan struct{})

			informer := informerFactory.Config().V1().ClusterVersions().Informer()
			go informerFactory.Start(stopCh)
			cache.WaitForCacheSync(stopCh, informer.HasSynced)

			informerFactory.Config().V1().ClusterVersions().Informer().GetStore().Add(tc.existingCV)

			r := &ClusterVersionReconciler{
				Client: fakeClient,
				Lister: informerFactory.Config().V1().ClusterVersions().Lister(),
				Log:    ctrl.Log.WithName("TestReconcile"),
			}

			request := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name: tc.existingCV.Name,
				},
			}

			_, err := r.Reconcile(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}

			cv, err := r.Client.ConfigV1().ClusterVersions().Get(context.Background(), tc.existingCV.Name, metav1.GetOptions{})
			if err != nil {
				t.Fatal(err)
			}

			if cv.Spec.DesiredUpdate != nil {
				t.Errorf("got: %v, expected: nil", cv.Spec.DesiredUpdate)
			}

			if cv.Spec.Channel != "" {
				t.Errorf("got: %v, expected: ''", cv.Spec.Channel)
			}

			if cv.Spec.Upstream != "" {
				t.Errorf("got: %v, expected: ''", cv.Spec.Upstream)
			}

			close(stopCh)
		})
	}
}
