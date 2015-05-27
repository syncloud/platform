#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp
NAME=python
VERSION=2.7.10
ROOT=/opt/syncloud-platform
PREFIX=${ROOT}/${NAME}

echo "building ${NAME}"

apt-get -y install build-essential

rm -rf build
mkdir -p build
cd build

wget https://www.python.org/ftp/python/${VERSION}/Python-${VERSION}.tgz
tar xzf ${NAME}-${VERSION}.tgz
cd ${NAME}-${VERSION}
./configure --prefix=${PREFIX}
make
rm -rf ${PREFIX}
make install
cd ..

tar cpzf ${NAME}.tar.gz -C ${ROOT} ${NAME}

cd ..