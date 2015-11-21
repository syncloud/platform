#!/bin/bash

DEVICE="/dev/mmcblk0"

PARTED_LINES=$(parted -sm ${DEVICE} unit B print | wc -l)
PARTITION_NUM=$(expr ${PARTED_LINES} - 2)
PARTITION="${DEVICE}p${PARTITION_NUM}"

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

resize2fs ${PARTITION}