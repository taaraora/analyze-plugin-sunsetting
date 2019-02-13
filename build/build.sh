#!/usr/bin/env bash

# Exit script when command fails
set -o errexit
# Exit script when it tries to use undeclared variables
set -o nounset
# if any of the commands in pipeline fails, script will exit
set -o pipefail


GO_FLAGS=${GO_FLAGS:-"-tags netgo"}    # Extra go flags to use in the build.
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
PLUGIN_NAME="analyze-plugin-sunsetting"
REPO_PATH="github.com/supergiant/analyze"

version=$( git describe --tags --dirty --abbrev=14 | sed -E 's/-([0-9]+)-g/.\1+/' )
revision=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )
go_version=$( go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/' )


ldflags="
  -X ${REPO_PATH}/version.Version=${version}
  -X ${REPO_PATH}/version.Revision=${revision}
  -X ${REPO_PATH}/version.Branch=${branch}
  -X ${REPO_PATH}/version.BuildDate=${BUILD_DATE}
  -X ${REPO_PATH}/version.GoVersion=${go_version}"

echo "Building $PLUGIN_NAME with -ldflags $ldflags"


GOBIN=$PWD go build ${GO_FLAGS} -ldflags "${ldflags}" "${REPO_PATH}"

exit 0
