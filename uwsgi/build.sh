#!/bin/bash


DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

export TMPDIR=/tmp
export TMP=/tmp

NAME=uwsgi
VERSION=2.0.10

echo "building ${NAME}"

apt-get -y install build-essential python-dev

rm -rf build
mkdir -p build
cd build

wget http://projects.unbit.it/downloads/${NAME}-${VERSION}.tar.gz
tar xzf ${NAME}-${VERSION}.tar.gz
cd ${NAME}-${VERSION}

python uwsgiconfig.py --build

cd ../..

mkdir -p install/${NAME}/bin
cp build/${NAME}-${VERSION}/uwsgi install/${NAME}/bin/

rm -rf ${NAME}.tar.gz
tar cpzf ${NAME}.tar.gz -C install ${NAME}