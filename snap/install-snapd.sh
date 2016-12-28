#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 version"
    exit 1
fi

VERSION=$1

tar xzvf snapd-${VERSION}-amd64.tar.gz
systemctl stop snapd.service snapd.socket || true
systemctl disable snapd.service snapd.socket || true

rm -rf /var/lib/snapd
mkdir /var/lib/snapd

rm -rf /usr/lib/snapd
mkdir -p /usr/lib/snapd
cp snapd/bin/snapd /usr/lib/snapd/snapd
cp snapd/bin/snap-exec /usr/lib/snapd/snap-exec
cp snapd/bin/snap-confine /usr/lib/snapd/snap-confine
cp snapd/bin/snap-discard-ns /usr/lib/snapd/snap-discard-ns
cp snapd/bin/snap /usr/bin/snap
cp snapd/bin/snapctl /usr/bin/snapctl

cp snapd/conf/snapd.service /lib/systemd/system/
cp snapd/conf/snapd.socket /lib/systemd/system/


systemctl enable snapd.service
systemctl enable snapd.socket
systemctl start snapd.service snapd.socket

snap install hello-world

snap --version

#TESTSLIB=${DIR}/snapd/scripts
#. ${DIR}/snapd/scripts/prepare.sh
#update_core_snap_with_snap_exec_snapctl

