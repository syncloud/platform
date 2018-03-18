#!/bin/bash

BOOT_PARTITION_INFO=$(lsblk -pP -o PKNAME,NAME,MOUNTPOINT | grep 'MOUNTPOINT="/"')
DEVICE=$(echo ${BOOT_PARTITION_INFO} | cut -d' ' -f1 | cut -d'=' -f2 | tr -d '"')
PARTITION=$(echo ${BOOT_PARTITION_INFO} | cut -d' ' -f2 | cut -d'=' -f2 | tr -d '"')
PARTITION_NUM=2

DEVICE_SIZE_BYTES=$(parted -sm ${DEVICE} unit B print | grep -oP "${DEVICE}:\K[0-9]*(?=B)")
PART_START_BYTES=$(parted -sm ${DEVICE} unit B print | grep -oP "^${PARTITION_NUM}:\K[0-9]*(?=B)")
PART_START_SECTORS=$(expr ${PART_START_BYTES} / 512)
PART_END_SECTORS=$(expr ${DEVICE_SIZE_BYTES} / 512 - 1)

echo "
p
d
${PARTITION_NUM}
p
n
p
${PARTITION_NUM}
${PART_START_SECTORS}
${PART_END_SECTORS}
p
w
q
" | fdisk ${DEVICE}

partprobe

resize2fs ${PARTITION}