#!/bin/bash

# Exit script when command fails
set -o errexit
# Exit script when it tries to use undeclared variables
set -o nounset
# if any of the commands in pipeline fails, script will exit
set -o pipefail

export TAG=${TRAVIS_BRANCH:-unstable}


### Main
echo "Tag Name: ${TAG}"
# If a tag is pushed, tests are run, the docker container is built and pushed to
# dockerhub, and then a release is pushed to the releases page.
if [[ "$TRAVIS_TAG" =~ ^v[0-9]. ]]; then
	echo "release"
	# run tests
	./run_tests.sh
	# Build Docker container
	./docker_build.sh
	# Push to Dockerhub
	./docker_push.sh
	# Push to releases page
	./push_release.sh
# on an unstable branch, tests are run and the docker container is built and pushed.
elif [[ "$TRAVIS_BRANCH" == *release-* ]]; then
	echo "unstable branch"
	export TAG="${TAG}-unstable"
	echo "Tag Name: ${TAG}"
	# run tests
	./run_tests.sh
	# Build docker container
	./docker_build.sh
	# Push to Dockerhub
	./docker_push.sh
# if a push to master happens, tests are only run
elif [[ "$TRAVIS_BRANCH" == "master" ]]; then
	echo "master branch - test will only be run"
	echo "Tag Name: ${TAG}"
	# run tests
	./run_tests.sh
else
# any other branch is considered a testing branch and will only run tests and build the container.
	echo "testing branch - run tests and docker build"
	export TAG="${TAG}-testing"
	echo "Tag Name: ${TAG}"
	# run tests
	./run_tests.sh
	# Build docker container
	./docker_build.sh
fi
