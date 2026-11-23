#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 version"
    exit 1
fi

NAME=platform
ARCH=$(uname -m)
VERSION=$1
cd ${DIR}/build

BUILD_DIR=${DIR}/build/snap

apt update -o Acquire::Retries=5
apt install -y wget squashfs-tools dpkg-dev

# wget treats DNS resolution failures as fatal and never retries them
download() {
    local url=$1
    local out=$2
    local attempt
    for attempt in $(seq 1 10); do
        if wget --tries=3 --timeout=60 --retry-connrefused --retry-on-http-error=503,429 --progress=dot:giga -O "${out}" "${url}"; then
            return 0
        fi
        echo "download attempt ${attempt} failed for ${url}, retrying in $((attempt * 5))s"
        sleep $((attempt * 5))
    done
    echo "download failed after 10 attempts: ${url}"
    return 1
}

cp -r ${DIR}/bin ${BUILD_DIR}
cp -r ${DIR}/config ${BUILD_DIR}

echo "deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian testing main" > /etc/apt/sources.list.d/testing.list
apt update -o Acquire::Retries=5
apt-get download -o Acquire::Retries=5 ca-certificates/testing
dpkg -x ca-certificates_*_all.deb .
mv usr/share/ca-certificates/mozilla ${BUILD_DIR}/certs

download https://github.com/syncloud/3rdparty/releases/download/openldap/openldap-${ARCH}.tar.gz openldap-${ARCH}.tar.gz
tar xf openldap-${ARCH}.tar.gz
mv openldap ${BUILD_DIR}

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
echo ${VERSION} > $DIR/version
mksquashfs ${BUILD_DIR} ${DIR}/${PACKAGE} -noappend -comp xz -no-xattrs -all-root
mkdir ${DIR}/artifact
cp ${DIR}/${PACKAGE} ${DIR}/artifact

