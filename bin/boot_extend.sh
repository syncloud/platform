#!/bin/bash -xe

BOOT_PARTITION_INFO=$(lsblk -pP -o PKNAME,NAME,MOUNTPOINT | grep 'MOUNTPOINT="/"')
DEVICE=$(echo ${BOOT_PARTITION_INFO} | cut -d' ' -f1 | cut -d'=' -f2 | tr -d '"')
PARTITION=$(echo ${BOOT_PARTITION_INFO} | cut -d' ' -f2 | cut -d'=' -f2 | tr -d '"')
PARTITION_NUM=2

DEVICE_SIZE_BYTES=$(parted -sm ${DEVICE} unit B print | grep "^${DEVICE}:" | cut -d':' -f2 | cut -d'B' -f1)
PART_START_BYTES=$(parted -sm ${DEVICE} unit B print | grep "^${PARTITION_NUM}:" | cut -d':' -f2 | cut -d'B' -f1)
PART_END_BYTES=$(parted -sm ${DEVICE} unit B print | grep "^${PARTITION_NUM}:" | cut -d':' -f3 | cut -d'B' -f1)
PART_START_SECTORS=$(expr ${PART_START_BYTES} / 512)
PART_END_SECTORS=$(expr ${DEVICE_SIZE_BYTES} / 512 - 1)
UNUSED_BYTES=$(( $DEVICE_SIZE_BYTES - $PART_END_BYTES ))
MIN_FREE_SPACE_LIMIT_BYTES=100000
if [[ $UNUSED_BYTES -lt $MIN_FREE_SPACE_LIMIT_BYTES ]]; then
  echo "unused space is: ${UNUSED_BYTES}b is less then min free space limit (${MIN_FREE_SPACE_LIMIT_BYTES}b), not extending"
  exit 0
fi

if parted -sm ${DEVICE} unit B print | grep "^3:"; then
  echo "3 or more partitions are not supported"
  exit 0
fi

if parted -sm ${DEVICE} unit B print | grep "btrfs"; then
  echo "btrfs not supported"
  exit 0
fi

PTTYPE=$(fdisk -l ${DEVICE} | grep "Disklabel type:" | awk '{ print $3 }')
if [[ $PTTYPE == "gpt" ]]; then
  GPT_BACKUP_HEADER_SIZE=33
  PART_END_SECTORS=$(expr ${PART_END_SECTORS} - ${GPT_BACKUP_HEADER_SIZE})

echo "
p
d
${PARTITION_NUM}
p
n
${PARTITION_NUM}
${PART_START_SECTORS}
${PART_END_SECTORS}
p
w
" | fdisk ${DEVICE}

else

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

fi

partprobe

resize2fs ${PARTITION}
