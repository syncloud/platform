#!/usr/bin/env bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"
BUILD_DIR=${DIR}/build/${NAME}

${BUILD_DIR}/python/bin/pip install -r ${DIR}/dev_requirements.txt

# We need to run tests on platform python as it has some libraries like openssl
mv ${BUILD_DIR}/python/bin/py.test ${BUILD_DIR}/python/bin/py.test_runner
cat <<'EOF' > ${BUILD_DIR}/python/bin/py.test
#!/bin/bash
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
export LD_LIBRARY_PATH=${DIR}/lib
${DIR}/bin/python.bin ${DIR}/bin/py.test_runner "$@"
EOF
chmod +x ${BUILD_DIR}/python/bin/py.test
cd ${DIR}/src
${BUILD_DIR}/python/bin/py.test test