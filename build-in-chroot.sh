#!/bin/bash -xe

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

if [[ -z "$1" || -z "$2" ]]; then
    echo "usage $0 app_arch app_version"
    exit 1
fi

ARCH=$1
ROOTFS=rootfs

function cleanup {
    mount | grep rootfs
    mount | grep rootfs | awk '{print "umounting "$1; system("umount "$3)}'
    mount | grep rootfs
    
    lsof 2>&1 | grep rootfs
    
    rm -rf ${ROOTFS}
    rm -rf tmp
}

cleanup || true

mkdir tmp
cp -r * tmp/ || true

if [ ! -f rootfs-${ARCH}.tar.gz ]; then
  wget http://build.syncloud.org:8111/guestAuth/repository/download/debian_rootfs_syncloud_${ARCH}/lastSuccessful/rootfs.tar.gz\
  -O rootfs-${ARCH}.tar.gz --progress dot:giga
else
  echo "skipping rootfs"
fi

echo "extracting rootfs"
mkdir -p ${ROOTFS}
tar xzf rootfs-${ARCH}.tar.gz -C ${ROOTFS}

chroot ${ROOTFS} /bin/bash -c "mount -t devpts devpts /dev/pts"
chroot ${ROOTFS} /bin/bash -c "mount -t proc proc /proc"

mkdir ${ROOTFS}/temp
export TEMP=/temp
export TMP=/temp
export TMPDIR=/temp

mkdir ${ROOTFS}/build
cp -r tmp/* ${ROOTFS}/build/ || true

if [ -f ${ROOTFS}/build/deps.sh ]; then
    chroot ${ROOTFS} /build/deps.sh
fi

chroot ${ROOTFS} /build/build.sh $@

cp ${ROOTFS}/build/*.tar.gz .

cleanup || true