#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=smartmontools
VERSION=7.5
BUILD=${DIR}/build
PREFIX=${BUILD}/${NAME}

apt update
apt -y install build-essential wget

rm -rf ${BUILD}
mkdir ${BUILD}
cd ${BUILD}

for i in $(seq 1 10); do
    if wget --tries=3 --timeout=60 https://github.com/smartmontools/smartmontools/releases/download/RELEASE_${VERSION//./_}/smartmontools-${VERSION}.tar.gz; then
        break
    fi
    if [ "$i" = "10" ]; then
        echo "smartmontools download failed"
        exit 1
    fi
    echo "retry smartmontools download"
    sleep 10
done
tar xf smartmontools-${VERSION}.tar.gz
cd smartmontools-${VERSION}
./configure --without-systemdsystemunitdir --without-libsystemd --without-selinux --without-libcap-ng --without-nvme-devicescan
make smartctl
mkdir -p ${PREFIX}/bin ${PREFIX}/lib
cp smartctl ${PREFIX}/bin
cp ${DIR}/smartctl.sh ${PREFIX}/bin
ldd smartctl | grep "=> /" | awk '{print $3}' | xargs -I{} cp --remove-destination {} ${PREFIX}/lib
cp --remove-destination $(readlink -f /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/ld-linux*.so*) ${PREFIX}/lib/ld.so
ls -la ${PREFIX}/lib
${PREFIX}/lib/ld.so --library-path ${PREFIX}/lib --list ${PREFIX}/bin/smartctl

BUILD_DIR=${DIR}/../build/snap
mkdir -p ${BUILD_DIR}
rm -rf ${BUILD_DIR}/${NAME}
mv ${PREFIX} ${BUILD_DIR}/${NAME}
