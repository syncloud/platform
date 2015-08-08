#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

ACTION=$1

TMP_DISK=/tmp/disk
SILENCE=" > /dev/null"
function add {
    dd if=/dev/zero bs=1M count=10 of=${TMP_DISK} status=none
    echo "
p
n
p



w
q
" | fdisk ${TMP_DISK} > /dev/null 2>&1

    kpartx -a ${TMP_DISK}
    LOOP=$(kpartx -l ${TMP_DISK} | head -1 | cut -d ' ' -f1 | cut -c1-5)
    mkfs.ext4 -q /dev/mapper/${LOOP}p1
    echo "/dev/mapper/${LOOP}p1"
}

function remove {
    LOOP=$(kpartx -l ${TMP_DISK} | head -1 | cut -d ' ' -f1 | cut -c1-5)
    umount /dev/mapper/${LOOP}p1 > /dev/null 2>&1
    kpartx -d ${TMP_DISK} > /dev/null
    rm ${TMP_DISK}
}

case "$1" in
        add)
            remove
            add
            ;;

        remove)
            remove
            ;;

        *)
            echo $"Usage: $0 {add|remove}"
            exit 1

esac

