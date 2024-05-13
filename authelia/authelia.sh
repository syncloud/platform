#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
${DIR}/lib/ld-musl-*.so* --library-path ${DIR}/lib ${DIR}/authelia "$@"
