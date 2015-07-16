#!/bin/bash

ARCH=$(dpkg-architecture -qDEB_HOST_GNU_CPU)
if [ ! -z "$1" ]; then
    ARCH=$1
fi

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