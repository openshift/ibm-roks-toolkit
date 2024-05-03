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
RELEASE_BRANCH="release-${RELEASE}"

if ! git rev-parse --verify "${RELEASE_BRANCH}" &> /dev/null; then
  echo "Skipping image release for master branch"
  exit
fi

if ! which oc &> /dev/null; then
  echo "ERROR: the oc command is required for this script."
  exit 1
fi
oc whoami

if ! oc get istag -n hypershift-toolkit ibm-roks-control-plane-operator:${RELEASE_TAG} &> /dev/null; then
  echo "Tagging control-plane-operator image with ${RELEASE_TAG}"
  oc tag registry.ci.openshift.org/hypershift-toolkit/ibm-roks-${RELEASE}:control-plane-operator \
        ibm-roks-control-plane-operator:${RELEASE_TAG} \
        -n hypershift-toolkit
else
  echo "control-plane-operator image tag ${RELEASE_TAG} already exists"
fi

if ! oc get istag -n hypershift-toolkit ibm-roks-toolkit:${RELEASE_TAG} &> /dev/null; then
  echo "Tagging ibm-roks-toolkit image with ${RELEASE_TAG}"
  oc tag registry.ci.openshift.org/hypershift-toolkit/ibm-roks-${RELEASE}:ibm-roks-toolkit \
        ibm-roks-toolkit:${RELEASE_TAG} \
        -n hypershift-toolkit
else
  echo "ibm-roks-toolkit image tag ${RELEASE_TAG} already exists"
fi

if ! oc get istag -n hypershift-toolkit ibm-roks-metrics:${RELEASE_TAG} &> /dev/null; then
  echo "Tagging ibm-roks-metrics image with ${RELEASE_TAG}"
  oc tag registry.ci.openshift.org/hypershift-toolkit/ibm-roks-${RELEASE}:metrics \
      ibm-roks-metrics:${RELEASE_TAG} \
      -n hypershift-toolkit
else
  echo "ibm-roks-metrics image tag ${RELEASE_TAG} already exists"
fi
