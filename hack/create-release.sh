#!/bin/bash
set -euo pipefail

REPODIR="$(dirname "$0")/.."

${REPODIR}/hack/create-git-release.sh
${REPODIR}/hack/tag-release-images.sh
