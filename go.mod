module github.com/openshift/ibm-roks-toolkit

go 1.15

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/openshift/api v0.0.0-20200922074107-8c5a4702c21a
	github.com/openshift/client-go v0.0.0-20200827190008-3062137373b5
	github.com/openshift/cluster-kube-apiserver-operator v0.0.0-alpha.0.0.20200901175228-fa89e5a96600
	github.com/openshift/cluster-openshift-apiserver-operator v0.0.0-alpha.0.0.20200827092600-713cb5655059
	github.com/openshift/cluster-openshift-controller-manager-operator v0.0.0-alpha.0.0.20200810151007-268aac45c717
	github.com/openshift/library-go v0.0.0-20200921120329-c803a7b7bb2c
	github.com/openshift/oc v0.0.0-alpha.0.0.20200930115932-6bfa7f97c8d4
	github.com/openshift/openshift-controller-manager v0.0.0-alpha.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/apiserver v0.19.0
	k8s.io/cli-runtime v0.19.0
	k8s.io/client-go v0.19.0
	k8s.io/component-base v0.19.0
	k8s.io/klog/v2 v2.3.0
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/yaml v1.2.0
)

replace (
	bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20190414153302-2ae31c8b6b30
	github.com/apcera/gssapi => github.com/openshift/gssapi v0.0.0-20161010215902-5fb4217df13b
	github.com/containers/image => github.com/openshift/containers-image v0.0.0-20190130162827-4bc6d24282b1
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20191121165722-d1d5f6476656
	github.com/golang/glog => github.com/openshift/golang-glog v0.0.0-20190322123450-3c92600d7533
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.7
	github.com/jteeuwen/go-bindata => github.com/jteeuwen/go-bindata v3.0.8-0.20151023091102-a0ff2567cfb7+incompatible
	github.com/onsi/ginkgo => github.com/openshift/onsi-ginkgo v1.2.1-0.20190125161613-53ca7dc85f60
	k8s.io/api => k8s.io/api v0.19.0-rc.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0-rc.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0-rc.3
	k8s.io/apiserver => k8s.io/apiserver v0.19.0-rc.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.0-rc.3
	k8s.io/client-go => k8s.io/client-go v0.19.0-rc.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.0-rc.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.0-rc.3
	k8s.io/code-generator => k8s.io/code-generator v0.19.0-rc.3
	k8s.io/component-base => k8s.io/component-base v0.19.0-rc.3
	k8s.io/cri-api => k8s.io/cri-api v0.19.0-rc.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.0-rc.3
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.0-rc.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.0-rc.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.0-rc.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.0-rc.3
	k8s.io/kubectl => k8s.io/kubectl v0.19.0-rc.3
	k8s.io/kubelet => k8s.io/kubelet v0.19.0-rc.3
	k8s.io/kubernetes => github.com/openshift/kubernetes v1.19.0-rc.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.0-rc.3
	k8s.io/metrics => k8s.io/metrics v0.19.0-rc.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.0-rc.3
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.19.0-rc.3
	k8s.io/sample-controller => k8s.io/sample-controller v0.19.0-rc.3
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)
