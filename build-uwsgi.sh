#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARCH=$(uname -m)

ARCH=$(uname -m)
BUILD_DIR=${DIR}/build/platform/python
docker ps -a -q --filter ancestor=python:syncloud --format="{{.ID}}" | xargs docker stop | xargs docker rm || true
docker rmi python:syncloud || true
docker build -t python:syncloud .
docker run python:syncloud python --help
docker run python:syncloud uwsgi --help
docker create --name=python python:syncloud
mkdir -p ${BUILD_DIR}
cd ${BUILD_DIR}
docker export python -o python.tar
tar xf python.tar
rm -rf python.tar
cp ${DIR}/bin/python ${BUILD_DIR}/bin
cp ${DIR}/bin/pip ${BUILD_DIR}/bin
rm -rf ${BUILD_DIR}/usr/src

#apt update
#apt install -y wget build-essential libsasl2-dev libldap2-dev libssl-dev libjansson-dev
cd ${DIR}

#BUILD_DIR=${DIR}/build/platform
#PYTHON_DIR=${BUILD_DIR}/python
#export PATH=${PYTHON_DIR}/bin:$PATH

#wget -c --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/python3-${ARCH}.tar.gz
#tar xf python3-${ARCH}.tar.gz
#mv python3 ${BUILD_DIR}/python

#cd ${DIR}
#export UWSGI_PROFILE_OVERRIDE=ssl=false
#export CPPFLAGS=-I${PYTHON_DIR}/include
#export LDFLAGS=-L${PYTHON_DIR}/lib
#export LD_LIBRARY_PATH=${PYTHON_DIR}/lib
