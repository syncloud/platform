#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

BUILD_DIR=${DIR}/../build/snap
mkdir -p /snap/platform
ln -sfn ${BUILD_DIR} /snap/platform/current

SMARTCTL=/snap/platform/current/smartmontools/bin/smartctl.sh

${SMARTCTL} --version
${SMARTCTL} -j --version | grep -q '"json_format_version"'
! ${BUILD_DIR}/smartmontools/lib/ld.so --library-path ${BUILD_DIR}/smartmontools/lib --list ${BUILD_DIR}/smartmontools/bin/smartctl | grep "not found"
