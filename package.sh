#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 version"
    exit 1
fi

NAME=platform
ARCH=$(uname -m)
VERSION=$1
CA_CERTIFICATES_VERSION=20241223
cd ${DIR}/build

BUILD_DIR=${DIR}/build/snap

apt update
apt install -y wget squashfs-tools dpkg-dev

cp -r ${DIR}/bin ${BUILD_DIR}
cp -r ${DIR}/config ${BUILD_DIR}

wget http://ftp.us.debian.org/debian/pool/main/c/ca-certificates/ca-certificates_${CA_CERTIFICATES_VERSION}_all.deb
dpkg -x ca-certificates_${CA_CERTIFICATES_VERSION}_all.deb .
mv usr/share/ca-certificates/mozilla ${BUILD_DIR}/certs

wget --retry-on-http-error=503 --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/nginx/nginx-${ARCH}.tar.gz
tar xf nginx-${ARCH}.tar.gz
mv nginx ${BUILD_DIR}
wget --retry-on-http-error=503 --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/gptfdisk/gptfdisk-${ARCH}.tar.gz
tar xf gptfdisk-${ARCH}.tar.gz
mv gptfdisk ${BUILD_DIR}
wget --retry-on-http-error=503 --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/openldap/openldap-${ARCH}.tar.gz
tar xf openldap-${ARCH}.tar.gz
mv openldap ${BUILD_DIR}
wget --retry-on-http-error=503 --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/btrfs/btrfs-${ARCH}.tar.gz
tar xf btrfs-${ARCH}.tar.gz
mv btrfs ${BUILD_DIR}

cd ${DIR}/build

echo "snapping"
ARCH=$(dpkg-architecture -q DEB_HOST_ARCH)

cp -r ${DIR}/meta ${BUILD_DIR}
echo ${VERSION} >> ${BUILD_DIR}/meta/version
echo "version: $VERSION" >> ${BUILD_DIR}/meta/snap.yaml
echo "architectures:" >> ${BUILD_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${BUILD_DIR}/meta/snap.yaml

PACKAGE=${NAME}_${VERSION}_${ARCH}.snap
echo ${PACKAGE} > $DIR/package.name
mksquashfs ${BUILD_DIR} ${DIR}/${PACKAGE} -noappend -comp xz -no-xattrs -all-root
mkdir ${DIR}/artifact
cp ${DIR}/${PACKAGE} ${DIR}/artifact

