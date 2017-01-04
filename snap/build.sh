#!/bin/bash -x

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
NAME="platform"

if [[ -z "$1" || -z "$2" ]]; then
    echo "usage $0 app_arch app_version"
    exit 1
fi

ARCH=$1
VERSION=$2
cd ${DIR}/..
./build.sh $ARCH $VERSION

cd snap

rm -rf build
mkdir build
rm -rf *.snap
cp -r ../build/platform/* build/
cp -r meta build/
cp snapcraft.yaml build/meta/snap.yaml
echo "version: $VERSION" >> build/meta/snap.yaml
echo "arch: $ARCH" >> build/meta/snap.yaml

mksquashfs build/ syncloud-platform_$VERSION_$ARCH.snap -noappend -comp xz -no-xattrs -all-root

