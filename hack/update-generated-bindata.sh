#!/bin/bash
set -euo pipefail

REPO_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )/.."

OUTPUTFILE="${OUTPUTFILE:-${REPO_DIR}/pkg/assets/bindata.go}"

TMP_GOPATH="$(mktemp -d)"

ln -s ${REPO_DIR}/vendor "${TMP_GOPATH}/src"

pushd ${REPO_DIR} &> /dev/null
GO111MODULE=off GOPATH="${TMP_GOPATH}" go install "./vendor/github.com/jteeuwen/go-bindata/..."
popd &> /dev/null

"${TMP_GOPATH}/bin/go-bindata" -nocompress -nometadata -prefix "${REPO_DIR}/assets/" -pkg assets -o "${OUTPUTFILE}" "${REPO_DIR}/assets/..."

gofmt -s -w "${OUTPUTFILE}"