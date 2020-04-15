#!/bin/bash

set -euo pipefail

TMP_OUTPUT="$(mktemp)"

REPO_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )/.."

OUTPUTFILE="${TMP_OUTPUT}" "${REPO_DIR}/hack/update-generated-bindata.sh"

diff -Naup ${REPO_DIR}/pkg/assets/bindata.go "${TMP_OUTPUT}"
