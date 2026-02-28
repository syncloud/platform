#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
LD=$(find ${DIR}/lib -name 'ld-*.so*' -type f | head -1)
exec ${LD} --library-path ${DIR}/lib ${DIR}/authelia "$@"
