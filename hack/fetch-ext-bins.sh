#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Enable tracing in this script off by setting the TRACE variable in your
# environment to any value:
#
# $ TRACE=1 test.sh
TRACE=${TRACE:-""}
if [ -n "$TRACE" ]; then
  set -x
fi

k8s_version=1.29.0
goarch=amd64
goos="unknown"

if [[ "$OSTYPE" == "linux"* ]]; then
  goos="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  goos="darwin"
fi

if [[ "$goos" == "unknown" ]]; then
  echo "OS '$OSTYPE' not supported. Aborting." >&2
  exit 1
fi

# Turn colors in this script off by setting the NO_COLOR variable in your
# environment to any value:
#
# $ NO_COLOR=1 test.sh
NO_COLOR=${NO_COLOR:-""}
if [ -z "$NO_COLOR" ]; then
  header=$'\e[1;33m'
  reset=$'\e[0m'
else
  header=''
  reset=''
fi

function header_text {
  echo "$header$*$reset"
}

rc=0
tmp_root=/tmp

envtest_root_dir=$tmp_root/envtest
envtest_bin_dir=$tmp_root/controller-tools/envtest

# Skip fetching and untaring the tools by setting the SKIP_FETCH_TOOLS variable
# in your environment to any value:
#
# $ SKIP_FETCH_TOOLS=1 ./fetch_ext_bins.sh
#
# If you skip fetching tools, this script will use the tools already on your
# machine.
SKIP_FETCH_TOOLS=${SKIP_FETCH_TOOLS:-""}

# fetch k8s API gen tools and make it available under envtest_bin_dir.
function fetch_tools {
  if [ -n "$SKIP_FETCH_TOOLS" ]; then
    return 0
  fi

  header_text "fetching tools"
  envtest_archive_name="envtest-v$k8s_version-$goos-$goarch.tar.gz"
  envtest_download_url="https://github.com/kubernetes-sigs/controller-tools/releases/download/envtest-v$k8s_version/$envtest_archive_name"

  envtest_archive_path="$tmp_root/$envtest_archive_name"
  if [ ! -f "$envtest_archive_path" ]; then
    curl -fsL "${envtest_download_url}" -o "$envtest_archive_path"
  fi

  tar -zvxf "$envtest_archive_path" -C "$tmp_root/"
}

function setup_envs {
  header_text "setting up env vars"

  # Setup env vars
  export PATH=$envtest_bin_dir:$PATH
  export TEST_ASSET_KUBECTL=$envtest_bin_dir/kubectl
  export TEST_ASSET_KUBE_APISERVER=$envtest_bin_dir/kube-apiserver
  export TEST_ASSET_ETCD=$envtest_bin_dir/etcd
  export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=10m
}
