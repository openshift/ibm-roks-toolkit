# IBM-ROKS Toolkit

## Overview
The IBM-ROKS toolkit is a set of tools and files that enables running OpenShift 4.x on IBM Public Cloud in a hyperscale manner with many control planes hosted on a central management cluster. 
This tool was jointly developed by RedHat and IBM.

## Getting Started

### Install on standalone environment

* Run `make build` to build the binary
* Construct a "cluster.yaml" to define custom parameters for the cluster. Example found here: [cluster.yaml.example](https://github.com/openshift/ibm-roks-toolkit/blob/master/cluster.yaml.example)
* Construct a "pull-secret.txt" to provide authentication to pull from desired docker registries. Example found here: [pull-secret.txt.example](https://github.com/openshift/ibm-roks-toolkit/blob/master/pull-secret.txt.example)
* Construct and run the render command, with optional fields below: `./bin/ibm-roks render`
    - `output-dir`: Specify the directory where manifest files should be output (default ./manifests)
    - `config`: Specify the config file for this cluster (default ./cluster.yaml)
    - `pull-secret`: Specify the pull secret used to pull from desired docker registries (default ./pull-secret.txt)
    - `pki-dir`: Specify the directory where the input PKI files have been placed (default ./pki)
    - `include-secrets`: If true, PKI secrets will be included in rendered manifests (default false)
    - `include-etcd`: If true, Etcd manifests will be included in rendered manifests (default false)
    - `include-autoapprover`: If true, includes a simple autoapprover pod in manifests (default false)
    - `include-vpn`: If true, includes a VPN server, sidecar and client (default false)
    - `include-registry`: If true, includes a default registry config to deploy into the user cluster (default false)
* Apply all the generated resources to the cluster `kubectl apply -f output-dir/`

## Release Process

### Creating a new release

New releases for the toolkit are created via pull requests.

1. Run `hack/bump-release.sh`. This will increment the date in the [release-date](release-date) file.
1. Submit a pull request with this change.
1. Once the PR is merged, a post submit job will automatically be kicked off to publish the release.
   - You can track the status of the post release jobs for all branches [here](https://prow.ci.openshift.org/?repo=openshift%2Fibm-roks-toolkit&type=postsubmit).
1. (optional) Cherry-pick the change to other branches.

### Changing the base image

The base images specified in the Dockerfiles are for testing only.
When releases are created and images published, the base images are substituted.
Master branch substitution definitions can be found [here](https://github.com/openshift/release/blob/master/ci-operator/config/openshift/ibm-roks-toolkit/openshift-ibm-roks-toolkit-master.yaml#L1-L9), with the other release branches located in the same directory.
There is a golang base, which is used for the first stage of the Dockerfiles, and a roks-toolkit-base, which is used for the second stage.
The roks-toolkit-base image is further defined [here](https://github.com/openshift/release/blob/master/clusters/app.ci/supplemental-ci-images/ibm-roks-toolkit-base/ibm-roks-toolkit-base.yaml).
To bump the base image for all toolkit images, submit a PR to update the image in that file.
Once merged, follow the steps for [Creating a new release](#creating-a-new-release).
