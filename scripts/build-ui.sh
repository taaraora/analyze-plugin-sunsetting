#!/usr/bin/env bash

# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail

BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
PLUGIN_NAME="analyze-plugin-sunsetting"
PLUGIN_VERSION=$( git describe --tags --dirty --abbrev=14 | sed -E 's/-([0-9]+)-g/.\1+/' )

//blabala popate environment.ts

cd ./ui && npm run build:webcomponents
