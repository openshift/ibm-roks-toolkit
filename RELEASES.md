# IBM-ROKS Toolkit

## Overview
For every PR that is merged into this repository, a post-commit CI job checks whether a release should be created. If so,
it will create a new release by publishing release artifacts to the git repository and tagging images in the CI registry
with the release tag.

Releases are controlled by 2 files in the root of the repository:

- release - Contains the major.minor version of the ibm-roks-toolkit branch. This should remain constant and only change
  when a new release branch is created.
- release-date - The release date in yyyymmdd format. This is used to create the release tag which is in the format:
  `v[RELEASE]-[RELEASE_DATE]`

Releases will only be created for named release branches and not for the master branch.

## Creating new Release

In the master branch, run `hack/bump-release.sh` to update the release-date file to the current date. Create a PR that 
contains this update and merge it. This should not result in a new release because it's the master branch. 

To create releases for named release branches, simply cherry-pick the PR to the other release branch.
When each PR merges, a new release will be created for that branch.
