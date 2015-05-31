#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=avahi
VERSION=0.6.31
ROOT=/opt/app/platform
PREFIX=${ROOT}/${NAME}

echo "building ${NAME}"

apt-get -y install build-essential intltool libgdbm3 libdaemon-dev libdbus-1-dev libdaemon0 dpkg-dev libcap-dev

cd build

wget http://avahi.org/download/${NAME}-${VERSION}.tar.gz
tar xzf ${NAME}-${VERSION}.tar.gz
cd ${NAME}-${VERSION}
./configure --prefix=${PREFIX} \
    --disable-qt3 \
    --disable-qt4 \
    --disable-gtk3 \
    --disable-gtk \
    --disable-mono \
    --disable-monodoc \
    --disable-python \
    --sysconfdir=${ROOT}/config/avahi \
    --enable-compat-libdns_sd \
    --with-systemdsystemunitdir=${PREFIX}/systemd
make
make install

cp --remove-destination /usr/lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/libdaemon.so* ${PREFIX}/lib

cd ../..