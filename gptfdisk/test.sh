#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

BUILD_DIR=${DIR}/../build/snap
mkdir -p /snap/platform
ln -sf ${BUILD_DIR} /snap/platform/current

SGDISK=/snap/platform/current/gptfdisk/bin/sgdisk.sh

${SGDISK} --help

# Test actual partitioning on a file-backed disk image
DISK=/tmp/test-disk.img
dd if=/dev/zero of=${DISK} bs=1M count=20

dd if=/dev/zero of=${DISK} bs=512 count=1 conv=notrunc
${SGDISK} -o ${DISK}
${SGDISK} -n 1 ${DISK}
${SGDISK} -p ${DISK}

rm -f ${DISK}
