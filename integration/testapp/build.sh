#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
rm -rf ${DIR}/*.snap
SNAP_DIR=${DIR}/build
mkdir ${SNAP_DIR}
cp -r ${DIR}/meta ${SNAP_DIR}/
echo "version: $VERSION" >> ${SNAP_DIR}/meta/snap.yaml
echo "architectures:" >> ${SNAP_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${SNAP_DIR}/meta/snap.yaml

mksquashfs ${SNAP_DIR} ${DIR}/testapp_1_${ARCH}.snap -noappend -comp xz -no-xattrs -all-root
cp ${DIR}/*.snap ${DIR}/artifact
