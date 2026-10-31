#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILD_DIR=${DIR}/build

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
mkdir -p ${BUILD_DIR}/bin ${BUILD_DIR}/meta
cp -r ${DIR}/bin/* ${BUILD_DIR}/bin/
cp ${DIR}/meta/snap.yaml ${BUILD_DIR}/meta/snap.yaml
echo "architectures:" >> ${BUILD_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${BUILD_DIR}/meta/snap.yaml

mksquashfs ${BUILD_DIR} ${DIR}/streamapp.snap -noappend -comp xz -no-xattrs -all-root
cp ${DIR}/*.snap ${DIR}/../../artifact
rm -rf ${BUILD_DIR}
