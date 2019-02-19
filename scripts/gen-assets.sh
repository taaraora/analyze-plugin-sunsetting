#!/usr/bin/env bash

# Exit script when command fails
set -o errexit
# Exit script when it tries to use undeclared variables
set -o nounset
# if any of the commands in pipeline fails, script will exit
set -o pipefail

# connect
cp  -R ./ui/dist/check ./asset
cp  -R ./ui/dist/settings ./asset

cd ./asset && go generate
