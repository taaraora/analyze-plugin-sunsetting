#!/usr/bin/env bash

# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail


BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
PLUGIN_NAME="analyze-plugin-sunsetting"
REPO_PATH="github.com/supergiant/analyze-plugin-sunsetting"


settings_component_entry_point="/settings/analyze-plugin-sunsetting-settings-main.js"
check_component_entry_point="/check/analyze-plugin-sunsetting-check-main.js"
version=$( git describe --tags --dirty --abbrev=14 | sed -E 's/-([0-9]+)-g/.\1+/' )
revision=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )
go_version=$( go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/' )


ldflags="
  -X ${REPO_PATH}/info.SettingsComponentEntryPoint=${settings_component_entry_point}
  -X ${REPO_PATH}/info.CheckComponentEntryPoint=${check_component_entry_point}
  -X ${REPO_PATH}/info.Version=${version}
  -X ${REPO_PATH}/info.Revision=${revision}
  -X ${REPO_PATH}/info.Branch=${branch}
  -X ${REPO_PATH}/info.BuildDate=${BUILD_DATE}
  -X ${REPO_PATH}/info.GoVersion=${go_version}"

echo "Building $PLUGIN_NAME with -ldflags $ldflags"


GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o ./dist/analyze-sunsetting -a -ldflags "${ldflags}" ./cmd/analyze-sunsetting
