#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILD_DIR=${DIR}/build

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
cp -r ${DIR}/bin/* ${BUILD_DIR}/bin/
cp ${DIR}/meta/snap.yaml ${BUILD_DIR}/meta/snap.yaml
cp -r ${DIR}/config ${BUILD_DIR}
cp -r ${DIR}/www ${BUILD_DIR}
cp -r ${DIR}/../../build/snap/nginx ${BUILD_DIR}
echo "architectures:" >> ${BUILD_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${BUILD_DIR}/meta/snap.yaml

mksquashfs ${BUILD_DIR} ${DIR}/testapp.snap -noappend -comp xz -no-xattrs -all-root
cp ${DIR}/*.snap ${DIR}/../../artifact
rm -rf ${BUILD_DIR}
