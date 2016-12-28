#!/bin/bash -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"

if [[ -z "$1" || -z "$2" ]]; then
    echo "usage $0 app_arch app_version"
    exit 1
fi

cd ${DIR}/..
./build.sh $1 $2

cd snap
snapcraft clean
rm -rf *.snap
snapcraft prime
cp -r meta prime/
snapcraft snap

