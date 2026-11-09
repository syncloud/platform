#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=lsblk
VERSION=$1
BUILD=${DIR}/build
PREFIX=${BUILD}/${NAME}

apt update
apt -y install build-essential pkg-config wget xz-utils libudev-dev

rm -rf ${BUILD}
mkdir ${BUILD}
cd ${BUILD}

for i in $(seq 1 10); do
    if wget --tries=3 --timeout=60 https://mirrors.edge.kernel.org/pub/linux/utils/util-linux/v${VERSION%.*}/util-linux-${VERSION}.tar.xz; then
        break
    fi
    if [ "$i" = "10" ]; then
        echo "util-linux download failed"
        exit 1
    fi
    echo "retry util-linux download"
    sleep 10
done
tar xf util-linux-${VERSION}.tar.xz
cd util-linux-${VERSION}
./configure --disable-all-programs --enable-lsblk --enable-libblkid --enable-libmount --enable-libsmartcols --with-udev --disable-shared --enable-static --without-python --without-systemd --disable-nls --disable-bash-completion --disable-asciidoc
make lsblk
mkdir -p ${PREFIX}/bin ${PREFIX}/lib
cp lsblk ${PREFIX}/bin
cp ${DIR}/lsblk.sh ${PREFIX}/bin
ldd lsblk | grep "=> /" | awk '{print $3}' | xargs -I{} cp --remove-destination {} ${PREFIX}/lib
cp --remove-destination $(readlink -f /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/ld-linux*.so*) ${PREFIX}/lib/ld.so
ls -la ${PREFIX}/lib
${PREFIX}/lib/ld.so --library-path ${PREFIX}/lib --list ${PREFIX}/bin/lsblk

BUILD_DIR=${DIR}/../build/snap
mkdir -p ${BUILD_DIR}
rm -rf ${BUILD_DIR}/${NAME}
mv ${PREFIX} ${BUILD_DIR}/${NAME}
