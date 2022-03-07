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
RELEASE_TAG="v${RELEASE}.0-$(cat "${REPODIR}/release-date")"

echo "The release tag is ${RELEASE_TAG}"

if ! which oc &> /dev/null; then
  echo "ERROR: the oc command is required for this script."
  exit 1
fi

if oc get istag -n hypershift-toolkit ibm-roks-control-plane-operator:${RELEASE_TAG} &> /dev/null; then
  echo "Removing control-plane-operator image tag ${RELEASE_TAG}"
  oc delete istag ibm-roks-control-plane-operator:${RELEASE_TAG} -n hypershift-toolkit
else
  echo "control-plane-operator image tag ${RELEASE_TAG} does not exist"
fi

if oc get istag -n hypershift-toolkit ibm-roks-toolkit:${RELEASE_TAG} &> /dev/null; then
  echo "Removing ibm-roks-toolkit image tag ${RELEASE_TAG}"
  oc delete istag ibm-roks-toolkit:${RELEASE_TAG} -n hypershift-toolkit
else
  echo "ibm-roks-toolkit image tag ${RELEASE_TAG} does not exist"
fi

if oc get istag -n hypershift-toolkit ibm-roks-metrics:${RELEASE_TAG} &> /dev/null; then
  echo "Removing ibm-roks-metrics image tag ${RELEASE_TAG}"
  oc delete istag ibm-roks-metrics:${RELEASE_TAG} -n hypershift-toolkit
else
  echo "ibm-roks-metrics image tag ${RELEASE_TAG} does not exist"
fi

