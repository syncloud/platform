#!/bin/bash -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"

if [[ -z "$1" || -z "$2" ]]; then
    echo "usage $0 app_arch app_version"
    exit 1
fi

VERSION=$2
cd ${DIR}/..
./build.sh $1 $VERSION

cd snap
sed 's/VERSION/$VERSION/g' -i snapcraft.yaml
snapcraft clean
rm -rf *.snap
snapcraft prime
cp -r meta prime/
snapcraft snap

