#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 device"
    exit 1
fi

DEVICE=$1
PARTITION=1

dd if=/dev/zero of=${DEVICE} bs=512 count=1 conv=notrunc
export LD_LIBRARY_PATH=${DIR}/gptfdisk/lib
${DIR}/gptfdisk/bin/sgdisk -o ${DEVICE}
${DIR}/gptfdisk/bin/sgdisk -n ${PARTITION} ${DEVICE}
${DIR}/gptfdisk/bin/sgdisk -p ${DEVICE}
PARTITION_DEVICE=$(lsblk -pl -o NAME,TYPE ${DEVICE} | grep part | awk '{print $1}')
mkfs.ext4 -F ${PARTITION_DEVICE}