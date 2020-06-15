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
echo "
o
p
n
p
${PARTITION}


p
w
q
" | ${DIR}/gptfdisk/bin/gdisk ${DEVICE}

mkfs.ext4 ${DEVICE}${PARTITION}