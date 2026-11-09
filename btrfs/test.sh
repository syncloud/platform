#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

BUILD_DIR=${DIR}/../build/snap
mkdir -p /snap/platform
ln -sfn ${BUILD_DIR} /snap/platform/current

BTRFS=/snap/platform/current/btrfs/bin/btrfs.sh
MKFS=/snap/platform/current/btrfs/bin/mkfs.sh

${BTRFS} version
${MKFS} --version

DISK=/tmp/btrfs-test.img
rm -f ${DISK}
truncate -s 200M ${DISK}
${MKFS} -f ${DISK}
${BTRFS} inspect-internal dump-super ${DISK} | grep -q "^fsid"
rm -f ${DISK}

! ${BUILD_DIR}/btrfs/lib/ld.so --library-path ${BUILD_DIR}/btrfs/lib --list ${BUILD_DIR}/btrfs/bin/btrfs.box | grep "not found"
test -f ${BUILD_DIR}/btrfs/lib/libgcc_s.so.1
