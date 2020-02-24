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
NODE_VERSION=10.15.1

cd ${DIR}

BUILD_DIR=${DIR}/build/${NAME}
GOROOT=${DIR}/go
PYTHON_DIR=${BUILD_DIR}/python
export PATH=${PYTHON_DIR}/bin:$GOROOT/bin:${DIR}/node/bin:$PATH
SNAP_DIR=${DIR}/build/snap
rm -rf build
mkdir -p ${BUILD_DIR}

cp -r ${DIR}/bin ${BUILD_DIR}

wget --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/nginx-${ARCH}.tar.gz
tar xf nginx-${ARCH}.tar.gz
mv nginx ${BUILD_DIR}
wget --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/uwsgi-${ARCH}.tar.gz
tar xf uwsgi-${ARCH}.tar.gz
mv uwsgi ${BUILD_DIR}
wget --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/openldap-${ARCH}.tar.gz
tar xf openldap-${ARCH}.tar.gz
mv openldap ${BUILD_DIR}
wget --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/openssl-${ARCH}.tar.gz
tar xf openssl-${ARCH}.tar.gz
mv openssl ${BUILD_DIR}
wget --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/python-${ARCH}.tar.gz
tar xf python-${ARCH}.tar.gz
mv python ${BUILD_DIR}

GO_ARCH=armv6l
NODE_ARCH=armv6l
if [[ ${ARCH} == "x86_64" ]]; then
    GO_ARCH=amd64
    NODE_ARCH=x64
fi

wget https://dl.google.com/go/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz --progress dot:giga
tar xf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

go version

wget https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${NODE_ARCH}.tar.gz \
    --progress dot:giga -O node.tar.gz
tar xzf node.tar.gz
mv node-v${NODE_VERSION}-linux-${NODE_ARCH} node

cd ${DIR}/www/public
npm install
npm run build

cd ${DIR}/backend
go test ./... -cover
CGO_ENABLED=0 go build -o ${BUILD_DIR}/bin/backend cmd/backend/main.go
CGO_ENABLED=0 go build -o ${BUILD_DIR}/bin/cli cmd/cli/main.go

cd ${DIR}

export CPPFLAGS=-I${PYTHON_DIR}/include
export LDFLAGS=-L${PYTHON_DIR}/lib
export LD_LIBRARY_PATH=${PYTHON_DIR}/lib

${PYTHON_DIR}/bin/pip install -r ${DIR}/requirements.txt

cd ${DIR}/src
rm -f version
echo ${VERSION} >> version
${PYTHON_DIR}/bin/python setup.py install
cd ..

cp -r ${DIR}/config ${BUILD_DIR}/config.templates
mkdir ${BUILD_DIR}/www
cp -r ${DIR}/www/internal ${BUILD_DIR}/www
cp -r ${DIR}/www/public/dist ${BUILD_DIR}/www/public

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
