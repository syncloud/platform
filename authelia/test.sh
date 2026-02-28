#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}
BUILD_DIR=${DIR}/../build/snap/authelia
ls -la ${BUILD_DIR}/
ls -laR ${BUILD_DIR}/lib/
ldd ${BUILD_DIR}/authelia
${BUILD_DIR}/authelia.sh -v