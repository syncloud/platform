#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
LD=$(ls ${DIR}/lib/ld-*.so* 2>/dev/null | head -1)
${LD} --library-path ${DIR}/lib ${DIR}/authelia "$@"
