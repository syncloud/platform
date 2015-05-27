#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
NAME=openldap
VERSION=1.8.0
ROOT=/opt/syncloud-platform
PREFIX=${ROOT}/${NAME}

echo "building ${NAME}"

apt-get -y install build-essential libdb-dev libdb++-dev

rm -rf build
mkdir -p build
cd build

wget http://www.openldap.org/software/download/OpenLDAP/openldap-release/${NAME}-${VERSION}.tgz
tar xzf ${NAME}-${VERSION}.tgz
cd ${NAME}-${VERSION}
./configure --prefix=${PREFIX}
make depend
make
rm -rf ${PREFIX}
make install
cd ..

tar cpzf ${NAME}.tar.gz -C ${ROOT} ${NAME}

cd ..