#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARCH=$(uname -m)
BUILD_DIR=${DIR}/build
mkdir ${BUILD_DIR}

rm -rf python3-${ARCH}.tar.gz
wget --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/python3-${ARCH}.tar.gz
tar xf python3-${ARCH}.tar.gz
mv python3 ${BUILD_DIR}/python
${BUILD_DIR}/python/bin/pip install -r ${DIR}/requirements.txt

ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
cp -r ${DIR}/meta ${BUILD_DIR}
cp -r ${DIR}/bin ${BUILD_DIR}
echo "architectures:" >> ${BUILD_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${BUILD_DIR}/meta/snap.yaml

mksquashfs ${BUILD_DIR} ${DIR}/testapp.snap -noappend -comp xz -no-xattrs -all-root
cp ${DIR}/*.snap ${DIR}/../../artifact