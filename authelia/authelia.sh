#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
LD=$(find ${DIR}/lib -name 'ld-*.so*' -type f -print -quit)
LIBS=$(echo ${DIR}/lib/*-linux-gnu*)
exec ${LD} --library-path $LIBS ${DIR}/authelia "$@"
