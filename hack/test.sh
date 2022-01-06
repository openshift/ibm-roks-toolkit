#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE}")/..

cd $REPO_ROOT && \
	source ./hack/fetch-ext-bins.sh && \
	fetch_tools && \
	setup_envs
	go test -v ./...

make test-render
