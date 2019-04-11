#!/bin/bash
# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail

# log into docker
echo "Log into docker"
docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
# push to docker
echo "Pushing to Docker"
docker push "$TRAVIS_REPO_SLUG":"$TAG"
