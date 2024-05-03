#!/bin/bash

set -exuo pipefail

REPODIR="$(dirname "$0")/.."
echo "$(date '+%Y%m%d')" > "${REPODIR}/release-date"
