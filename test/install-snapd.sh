#!/bin/bash -xe

VERSION=$(curl http://apps.syncloud.org/releases/stable/snapd2.version)
ARCH=$(dpkg --print-architecture)
SNAPD=snapd-${VERSION}-${ARCH}.tar.gz

cd /tmp
rm -rf "${SNAPD}"
rm -rf snapd
wget http://apps.syncloud.org/apps/"${SNAPD}" --progress=dot:giga
tar xzvf "${SNAPD}"
mkdir -p /var/lib/snapd/snaps
./snapd/install.sh
