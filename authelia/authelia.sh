#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
LIBS=$(echo ${DIR}/lib/*-linux-gnu*)
exec ${DIR}/lib/*-linux*/ld-linux*.so* --library-path $LIBS ${DIR}/authelia "$@"
