#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

NAME=dbus
VERSION=1.8.18
ROOT=/opt/app/platform
PREFIX=${ROOT}/avahi

echo "building ${NAME}"

cd build

wget http://dbus.freedesktop.org/releases/dbus/${NAME}-${VERSION}.tar.gz
tar xzf ${NAME}-${VERSION}.tar.gz
cd ${NAME}-${VERSION}
./configure --prefix=${PREFIX} \
    --sysconfdir=${ROOT}/config/avahi \
    --with-systemdsystemunitdir=${PREFIX}/systemd
make
make install
cd ../..