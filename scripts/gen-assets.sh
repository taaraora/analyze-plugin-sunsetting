#!/usr/bin/env bash

# Exit script when command fails
set -o errexit
# Exit script when it tries to use undeclared variables
set -o nounset
# if any of the commands in pipeline fails, script will exit
set -o pipefail


cd ./asset && go generate
