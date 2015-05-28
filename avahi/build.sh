#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
NAME=openldap
VERSION=0.6.31
ROOT=/opt/syncloud-platform
PREFIX=${ROOT}/${NAME}

echo "building ${NAME}"

apt-get -y install build-essential intltool libgtk2.0-dev qt4-qmake libqt4-dev libgdbm3 libdaemon-dev

rm -rf build
mkdir -p build
cd build

wget http://avahi.org/download/${NAME}-${VERSION}.tar.gz
tar xzf ${NAME}-${VERSION}.tar.gz
cd ${NAME}-${VERSION}
./configure --prefix=${PREFIX} --disable-qt3 --disable-gtk3 --disable-dbus --with-systemdsystemunitdir=${PREFIX}/systemd
make depend
make
rm -rf ${PREFIX}
make install
cd ..

tar cpzf ${NAME}.tar.gz -C ${ROOT} ${NAME}

cd ..