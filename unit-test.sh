#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"
BUILD_DIR=${DIR}/build/${NAME}

cd ${DIR}/src
${BUILD_DIR}/python/bin/pip install -r dev_requirements.txt
${BUILD_DIR}/python/bin/py.test test