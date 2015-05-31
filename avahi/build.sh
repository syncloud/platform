#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
NAME=avahi
VERSION=0.6.31
ROOT=/opt/app/platform
PREFIX=${ROOT}/${NAME}

echo "building ${NAME}"

apt-get -y install build-essential intltool libgdbm3 libdaemon-dev libdbus-1-dev libdaemon0 dpkg-dev

rm -rf build
mkdir -p build

rm -rf ${PREFIX}

./build-dbus.sh
./build-avahi.sh

rm -rf ${NAME}.tar.gz
tar cpzf ${NAME}.tar.gz -C ${ROOT} ${NAME}

cd ..