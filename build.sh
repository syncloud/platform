#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 version"
    exit 1
fi

NAME=$1
ARCH=$(uname -m)
VERSION=$2
GO_VERSION=1.11.5

cd ${DIR}

BUILD_DIR=${DIR}/build/${NAME}
GOROOT=${DIR}/go
PYTHON_DIR=${BUILD_DIR}/python
export PATH=${PYTHON_DIR}/bin:$GOROOT/bin:$PATH
SNAP_DIR=${DIR}/build/snap
rm -rf build
mkdir -p ${BUILD_DIR}

cp -r ${DIR}/bin ${BUILD_DIR}

DOWNLOAD_URL=http://artifact.syncloud.org/3rdparty
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/nginx-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/uwsgi-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/openldap-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/openssl-${ARCH}.tar.gz
coin --to ${BUILD_DIR} raw ${DOWNLOAD_URL}/python-${ARCH}.tar.gz

GO_ARCH=armv6l
if [[ ${ARCH} == "x86_64" ]]; then
    GO_ARCH=amd64
fi

wget https://dl.google.com/go/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz --progress dot:giga
tar xf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

go version
mkdir -p ${BUILD_DIR}/go/src/syncloud/platform
cp -r ${DIR}/backend/. ${BUILD_DIR}/go/src/syncloud/platform
export GOPATH=${BUILD_DIR}/go
cd ${BUILD_DIR}/go/src/syncloud/platform
go build
ls
ls ${BUILD_DIR}/go/bin
./backend
cp backend ${BUILD_DIR}/bin
rm -rf ${BUILD_DIR}/go

export CPPFLAGS=-I${PYTHON_DIR}/include
export LDFLAGS=-L${PYTHON_DIR}/lib
export LD_LIBRARY_PATH=${PYTHON_DIR}/lib

pip install -r ${DIR}/requirements.txt

cd ${DIR}/src
rm -f version
echo ${VERSION} >> version
${PYTHON_DIR}/bin/python setup.py install
cd ..

cp -r ${DIR}/config ${BUILD_DIR}/config.templates
cp -r ${DIR}/www ${BUILD_DIR}

mkdir ${BUILD_DIR}/META
echo ${NAME} >> ${BUILD_DIR}/META/app
echo ${VERSION} >> ${BUILD_DIR}/META/version

echo "snapping"
ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)
rm -rf ${DIR}/*.snap

mkdir ${SNAP_DIR}
cp -r ${BUILD_DIR}/* ${SNAP_DIR}/
cp -r ${DIR}/snap/meta ${SNAP_DIR}/
cp ${DIR}/snap/snap.yaml ${SNAP_DIR}/meta/snap.yaml
echo "version: $VERSION" >> ${SNAP_DIR}/meta/snap.yaml
echo "architectures:" >> ${SNAP_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${SNAP_DIR}/meta/snap.yaml
PACKAGE=${NAME}_${VERSION}_${ARCH}.snap
echo ${PACKAGE} > package.name

mksquashfs ${SNAP_DIR} ${DIR}/${PACKAGE} -noappend -comp xz -no-xattrs -all-root

