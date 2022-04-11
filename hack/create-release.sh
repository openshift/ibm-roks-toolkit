#!/bin/bash
set -euo pipefail

REPODIR="$(dirname "$0")/.."

${REPODIR}/hack/wait-for-images.sh
${REPODIR}/hack/create-git-release.sh
${REPODIR}/hack/tag-release-images.sh
