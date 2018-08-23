#!/bin/bash -e

if [[ -z "$1" ]]; then
    echo "usage $0 device"
    exit 1
fi

DEVICE=$1

dd if=/dev/zero of=${DEVICE} bs=512 count=1 conv=notrunc

echo "
o
p
n
p
1


p
w
q
" | fdisk ${DEVICE}

mkfs.ext4 ${DEVICE}