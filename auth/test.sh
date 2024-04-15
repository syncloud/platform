#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}
BUILD_DIR=${DIR}/../build/snap/auth
ldd ${BUILD_DIR}/authelia
${BUILD_DIR}/authelia.sh -v