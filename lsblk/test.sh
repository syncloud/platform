#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

BUILD_DIR=${DIR}/../build/snap
mkdir -p /snap/platform
ln -sfn ${BUILD_DIR} /snap/platform/current

LSBLK=/snap/platform/current/lsblk/bin/lsblk.sh

${LSBLK} --version
${LSBLK} -Jp
! ${BUILD_DIR}/lsblk/lib/ld.so --library-path ${BUILD_DIR}/lsblk/lib --list ${BUILD_DIR}/lsblk/bin/lsblk | grep "not found"
