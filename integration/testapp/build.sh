#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARCH=$(uname -m)
BUILD_DIR=${DIR}/build
mkdir ${BUILD_DIR}

cp -R ${DIR}/../build/platform/python ${BUILD_DIR}/python

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
cp -r ${DIR}/meta ${BUILD_DIR}
cp -r ${DIR}/bin ${BUILD_DIR}
echo "architectures:" >> ${BUILD_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${BUILD_DIR}/meta/snap.yaml

mksquashfs ${BUILD_DIR} ${DIR}/testapp.snap -noappend -comp xz -no-xattrs -all-root
cp ${DIR}/*.snap ${DIR}/../../artifact