#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
${DIR}/lib/ld-linux*.so* --library-path ${DIR}/lib ${DIR}/bin/sgdisk "$@"
