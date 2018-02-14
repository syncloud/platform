#!/bin/bash -e

SOURCE_DEVICE_NAME="/dev/mmcblk0"
TARGET_DEVICE_NAME="/dev/mmcblk1"

SOURCE_DEVICE_SIZE=$(blockdev --getsize64 ${SOURCE_DEVICE_NAME})
TARGET_DEVICE_SIZE=$(blockdev --getsize64 ${TARGET_DEVICE_NAME})
BYTES_TO_COPY=${SOURCE_DEVICE_SIZE}

echo "source device ($SOURCE_DEVICE_NAME) size: ${SOURCE_DEVICE_SIZE} bytes"
echo "target device ($TARGET_DEVICE_NAME) size: ${TARGET_DEVICE_SIZE} bytes"

if [[ ${TARGET_DEVICE_SIZE} -lt ${SOURCE_DEVICE_SIZE} ]]; then
    echo "target device size is less then target device"
    if [[ $1 == "-f" ]]; then
        BYTES_TO_COPY=${TARGET_DEVICE_SIZE}
    else
        echo "use -f to copy only bytes to fit on target (dangerous)"
        exit 1
    fi
fi

KIBI_BYTE=1024
KIBI_BYTES_TO_COPY=$(($BYTES_TO_COPY/$KIBI_BYTE))

dd if=${SOURCE_DEVICE_NAME} of=${TARGET_DEVICE_NAME} bs=${KIBI_BYTE} count=${KIBI_BYTES_TO_COPY}