#!/bin/bash

export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
export DEBCONF_FRONTEND=noninteractive
export DEBIAN_FRONTEND=noninteractive
export TMPDIR=/tmp
export TMP=/tmp

if [[ -z "$1" || -z "$2" ]]; then
    echo "usage $0 app_arch app_version"
    exit 1
fi

ARCH=$1

if [ ! -f "rootfs.tar.gz" ]; then
  wget http://build.syncloud.org:8111/guestAuth/repository/download/debian_rootfs_${ARCH}/lastSuccessful/rootfs.tar.gz\
  -O rootfs.tar.gz --progress dot:giga
else
    echo "rootfs.tar.gz is here"
fi

rm -rf /tmp/rootfs
mkdir /tmp/rootfs
tar xzf rootfs.tar.gz -C /tmp/rootfs
cp -r ./* /tmp/rootfs/root
chroot /tmp/rootfs root/build.sh $@
rm -rf platform*.tar.gz
mv /tmp/rootfs/root/platform*.tar.gz .
