#!/bin/bash

set -euo pipefail

REPODIR="$(dirname "$0")/.."
echo "$(date '+%Y%m%d')" > "${REPODIR}/release-date"
