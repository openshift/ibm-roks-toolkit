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
RELEASE_DATE="$(cat "${REPODIR}/release-date")"
GIT_RELEASE_TAG="v${RELEASE}.0+${RELEASE_DATE}"
RELEASE_BRANCH="release-${RELEASE}"

# if the branch corresponding to the release doesn't exist, skip because 
# it's likely the master branch
if ! git rev-parse --verify "${RELEASE_BRANCH}" &> /dev/null; then
  echo "Skipping git release for master branch"
  exit
fi

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "${RELEASE_BRANCH}" != "${CURRENT_BRANCH}" ]]; then
  echo "The current branch ${CURRENT_BRANCH} is not the expected release branch ${RELEASE_BRANCH}"
  exit 1
fi

if git rev-list ${GIT_RELEASE_TAG}.. &>/dev/null; then
  echo "The release ${GIT_RELEASE_TAG} already exists. Nothing to do"
  exit
fi

echo "Creating release ${GIT_RELEASE_TAG}"

if ! git remote | grep origin; then
  echo "Adding git remote"
  git remote add origin https://github.com/openshift/ibm-roks-toolkit.git
else
  echo "Setting git remote URL" 
  git remote set-url origin https://github.com/openshift/ibm-roks-toolkit.git
fi

git tag "${GIT_RELEASE_TAG}"
GORELEASER_CURRENT_TAG="${GIT_RELEASE_TAG}" goreleaser release --skip-publish --config "${REPODIR}/hack/release-config.yaml"

GITHUB_TOKEN="${GITHUB_TOKEN:-}"
if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "GITHUB_TOKEN is not present, will skip publishing release"
  exit
fi

hub release create \
  -a "./dist/ibm-roks-toolkit_${RELEASE}.0+${RELEASE_DATE}_linux_x86_64.tar.gz" \
  -a "./dist/checksums.txt" \
  -m "${GIT_RELEASE_TAG} ${RELEASE_DATE}" \
  -t "${RELEASE_BRANCH}" \
  "${GIT_RELEASE_TAG}"
