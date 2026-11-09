#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=btrfs
VERSION=$1
BUILD=${DIR}/build
PREFIX=${BUILD}/${NAME}

apt update
apt -y install build-essential autoconf automake pkg-config wget uuid-dev libblkid-dev zlib1g-dev liblzo2-dev libzstd-dev libudev-dev

rm -rf ${BUILD}
mkdir ${BUILD}
cd ${BUILD}

for i in $(seq 1 10); do
    if wget --tries=3 --timeout=60 -O btrfs-progs.tar.gz https://github.com/kdave/btrfs-progs/archive/refs/tags/v${VERSION}.tar.gz; then
        break
    fi
    if [ "$i" = "10" ]; then
        echo "btrfs-progs download failed"
        exit 1
    fi
    echo "retry btrfs-progs download"
    sleep 10
done
tar xf btrfs-progs.tar.gz
cd btrfs-progs-${VERSION}
./autogen.sh
./configure --disable-documentation --disable-convert --disable-python
make btrfs.box
mkdir -p ${PREFIX}/bin ${PREFIX}/lib
cp btrfs.box ${PREFIX}/bin/btrfs.box
cp btrfs.box ${PREFIX}/bin/mkfs.btrfs
cp ${DIR}/btrfs.sh ${DIR}/mkfs.sh ${PREFIX}/bin
ldd btrfs.box | grep "=> /" | awk '{print $3}' | xargs -I{} cp --remove-destination {} ${PREFIX}/lib
cp --remove-destination $(ldconfig -p | grep -m1 'libgcc_s.so.1 ' | awk '{print $NF}') ${PREFIX}/lib
cp --remove-destination $(readlink -f /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/ld-linux*.so*) ${PREFIX}/lib/ld.so
ls -la ${PREFIX}/lib
${PREFIX}/lib/ld.so --library-path ${PREFIX}/lib --list ${PREFIX}/bin/btrfs.box

BUILD_DIR=${DIR}/../build/snap
mkdir -p ${BUILD_DIR}
rm -rf ${BUILD_DIR}/${NAME}
mv ${PREFIX} ${BUILD_DIR}/${NAME}
