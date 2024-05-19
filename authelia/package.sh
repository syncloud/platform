#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}
BUILD_DIR=${DIR}/../build/snap/authelia
mkdir -p ${BUILD_DIR}
cp /app/authelia ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}
cp ${DIR}/authelia.sh ${BUILD_DIR}
