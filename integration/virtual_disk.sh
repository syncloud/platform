#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

apt-get install -y kpartx > /var/log/virtual_disk.log 2>&1

ACTION=$1
FS=$2

TMP_DISK=/tmp/disk
function add {
    dd if=/dev/zero bs=1M count=10 of=${TMP_DISK} status=none
    echo "
p
o
p
n
p



w
q
" | fdisk ${TMP_DISK} >> /var/log/virtual_disk.log 2>&1

    kpartx -as ${TMP_DISK}
    LOOP=$(kpartx -l ${TMP_DISK} | head -1 | cut -d ' ' -f1 | cut -c1-5)
    mkfs.${FS} /dev/mapper/${LOOP}p1 >> /var/log/virtual_disk.log 2>&1
    echo "/dev/mapper/${LOOP}p1"
}

function remove {
    LOOP=$(kpartx -l ${TMP_DISK} | head -1 | cut -d ' ' -f1 | cut -c1-5)
    umount /dev/mapper/${LOOP}p1 >> /var/log/virtual_disk.log 2>&1
    kpartx -d ${TMP_DISK} >> /var/log/virtual_disk.log 2>&1
    rm ${TMP_DISK} >> /var/log/virtual_disk.log 2>&1
}

case "$ACTION" in
        add)
            remove
            add
            ;;

        remove)
            remove
            ;;

        *)
            echo $"Usage: $0 {add|remove} fs"
            exit 1

esac

