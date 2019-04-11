#!/bin/bash
# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail

echo "Compressing assests"
tar -czpf /tmp/assest.gz "${TRAVIS_HOME}"/gopath/src/github.com/"${TRAVIS_REPO_SLUG}"/

# if a tag has alpha or beta in the name, it will be released as a pre-release.
# if a tag does not have alpha or beta, it is pushed as a full release.
case "${TAG}" in
	*alpha* )  echo "Releasing version: ${TAG}, as pre-release"
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "pre-release" --prerelease --debug "$TAG" /tmp/assest.gz;;
	*beta* )    echo "Releasing version: ${TAG}, as pre-release"
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "pre-release" --prerelease --debug "$TAG" /tmp/assest.gz;;
	*)echo "Releasing version: ${TAG}, as latest release."
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "latest release" --debug "$TAG" /tmp/assest.gz;;
esac
