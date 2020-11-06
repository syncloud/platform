#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARCH=$(uname -m)

cd ${DIR}

BUILD_DIR=${DIR}/build/platform
PYTHON_DIR=${BUILD_DIR}/python
export PATH=${PYTHON_DIR}/bin:$PATH

wget -c --progress=dot:giga https://github.com/syncloud/3rdparty/releases/download/1/python3-${ARCH}.tar.gz
tar xf python3-${ARCH}.tar.gz
mv python3 ${BUILD_DIR}/python

cd ${DIR}
export UWSGI_PROFILE_OVERRIDE=ssl=false
export CPPFLAGS=-I${PYTHON_DIR}/include
export LDFLAGS=-L${PYTHON_DIR}/lib
export LD_LIBRARY_PATH=${PYTHON_DIR}/lib

${PYTHON_DIR}/bin/pip install -r ${DIR}/requirements.txt
ldd ${PYTHON_DIR}/bin/uwsgi
cp --remove-destination /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/libz.so* ${PYTHON_DIR}/lib
cp --remove-destination /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/libuuid.so* ${PYTHON_DIR}/lib
cp --remove-destination /usr/lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)/libjansson.so* ${PYTHON_DIR}/lib
${PYTHON_DIR}/bin/uwsgi --help
