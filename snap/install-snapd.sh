#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 snapd_file"
    exit 1
fi

SNAPD_FILE=$1

apt-get install -y dpgk-dev

tar xzvf ${SNAPD_FILE}
systemctl stop snapd.service snapd.socket || true
systemctl disable snapd.service snapd.socket || true

rm -rf /var/lib/snapd
mkdir /var/lib/snapd

rm -rf /usr/lib/snapd
mkdir -p /usr/lib/snapd
cp snapd/bin/snapd /usr/lib/snapd
cp snapd/bin/snap-exec /usr/lib/snapd
cp snapd/bin/snap-confine /usr/lib/snapd
cp snapd/bin/snap-discard-ns /usr/lib/snapd
cp snapd/bin/snap /usr/bin
cp snapd/bin/snapctl /usr/bin
cp snapd/bin/mksquashfs /usr/bin
cp snapd/bin/unsquashfs /usr/bin
cp snapd/lib/* /lib/$(dpkg-architecture -q DEB_HOST_GNU_TYPE)

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

