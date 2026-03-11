#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

ARCH=$(uname -m)
NAME=gptfdisk
VERSION=1.0.10
BUILD=${DIR}/build
PREFIX=${BUILD}/${NAME}

apt update
apt -y install build-essential libpopt-dev uuid-dev wget

rm -rf ${BUILD}
mkdir ${BUILD}
cd ${BUILD}

wget http://www.rodsbooks.com/gdisk/gptfdisk-${VERSION}.tar.gz
tar xf gptfdisk-${VERSION}.tar.gz
cd gptfdisk-${VERSION}
make sgdisk
echo "=== ldd sgdisk ==="
ldd sgdisk
mkdir -p ${PREFIX}/bin
cp sgdisk ${PREFIX}/bin
cp ${DIR}/sgdisk.sh ${PREFIX}/bin
mkdir -p ${PREFIX}/lib
ldd sgdisk | grep "=> /" | awk '{print $3}' | xargs -I{} cp --remove-destination {} ${PREFIX}/lib
cp --remove-destination /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/ld-linux*.so* ${PREFIX}/lib
echo "=== bundled libs ==="
ls -la ${PREFIX}/lib/
echo "=== validate: ldd using bundled ld ==="
${PREFIX}/lib/ld-linux*.so* --library-path ${PREFIX}/lib --list ${PREFIX}/bin/sgdisk

BUILD_DIR=${DIR}/../build/snap
mkdir -p ${BUILD_DIR}
mv ${PREFIX} ${BUILD_DIR}/gptfdisk
