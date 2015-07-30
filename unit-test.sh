#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

PYTHON=${DIR}/build/platform/python/bin
export LD_LIBRARY_PATH=${DIR}/build/platform/python/lib

wget -O get-pip.py https://bootstrap.pypa.io/get-pip.py
${PYTHON}/python get-pip.py
rm get-pip.py

cd ${DIR}/src
${PYTHON}/python ${PYTHON}/pip2 install -U pytest
${PYTHON}/python ${PYTHON}/pip2 install -r dev_requirements.txt

${PYTHON}/py.test test