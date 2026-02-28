#!/bin/bash -ex
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
LD=$(find ${DIR}/lib -name 'ld-*.so*' -type f -print -quit)
exec ${LD} --library-path ${DIR}/lib ${DIR}/authelia "$@"
