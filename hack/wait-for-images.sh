#!/bin/bash
set -euo pipefail

REPODIR="$(dirname "$0")/.."

if [[ ! -f "${REPODIR}/release" ]]; then
  echo "Release name file (release) does not exist. Nothing to do"
fi
if [[ ! -f "${REPODIR}/release-date" ]]; then
  echo "Release date file (release-date) does not exist. Nothing to do"
fi
RELEASE="$(cat "${REPODIR}/release")"
RELEASE_BRANCH="release-${RELEASE}"

if ! which oc &> /dev/null; then
  echo "ERROR: the oc command is required for this script."
  exit 1
fi

if ! git rev-parse --verify "${RELEASE_BRANCH}" &> /dev/null; then
  echo "Skipping image release for master branch"
  exit
fi

CURRENT_COMMIT="$(git rev-parse "${RELEASE_BRANCH}")"

# Wait up to 45 minutes, otherwise exit
timeout=45

while [ $timeout -gt 0 ]; do
  # todo remove echo
  IMG=$(oc get istag ibm-roks-"${RELEASE}":metrics -n hypershift-toolkit -ojson)
  echo ${IMG}
  URI=$(oc get istag ibm-roks-"${RELEASE}":metrics -n hypershift-toolkit -o jsonpath='{.image.dockerImageManifests[0].digest}')
  podman pull "$URI"
  podman inspect "$URI"
  image_commit=$(podman inspect "$URI" | jq -r .[0].ContainerConfig.Labels.io\.openshift\.build\.commit\.id)
  if [[ $image_commit == "$CURRENT_COMMIT" ]]; then
    echo "Tag with expected commit found ${image_commit}"
    break
  fi
  echo "${timeout}: Waiting for image commit ${CURRENT_COMMIT}. Current image commit: ${image_commit}"
  sleep 60
  timeout=$(( $timeout - 1 ))
done

if [ $timeout -eq 0 ]; then
  echo "Timed out waiting for commit"
  exit 1
fi
