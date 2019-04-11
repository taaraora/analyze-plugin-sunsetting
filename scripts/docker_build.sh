#!/bin/bash
# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail

echo "$TRAVIS_REPO_SLUG":"$TAG"
# build the docker container
echo "Building Docker container"
docker build --tag build --target build .
docker build --tag "$TRAVIS_REPO_SLUG":"$TAG" .
