#!/bin/bash -e
#
# Copyright (c) 2013-2014 Robert Nelson <robertcnelson@gmail.com>
# Portions copyright (c) 2014 Charles Steinkuehler <charles@steinkuehler.net>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

#This script assumes, these packages are installed, as network may not be setup
#dosfstools initramfs-tools rsync u-boot-tools

if ! id | grep -q root; then
	echo "must be run as root"
	exit
fi

apt -y install rsync dosfstools

if ! lsblk /dev/mmcblk0p2 | grep "/$"; then
    echo "root fs must be /dev/mmcblk0p2"
    exit 1
fi

boot_drive=$(LC_ALL=C lsblk -l | grep "/boot/uboot" | awk '{print $1}')

if [ ${boot_drive} == "" ]; then
    mount /dev/mmcblk0p1 /boot/uboot
    mount -t tmpfs tmpfs /tmp
fi

source="/dev/mmcblk0"
destination="/dev/mmcblk1"

flush_cache () {
	sync
	blockdev --flushbufs ${destination}
}

# umount does not like device names without a valid /etc/mtab
# find the mount point from /proc/mounts
dev2dir () {
	grep -m 1 '^$1 ' /proc/mounts | while read LINE ; do set -- $LINE ; echo $2 ; done
}

write_failure () {
	echo "writing to [${destination}] failed..."

	[ -e /proc/$CYLON_PID ]  && kill $CYLON_PID > /dev/null 2>&1

	if [ -e /sys/class/leds/beaglebone\:green\:usr0/trigger ] ; then
		echo heartbeat > /sys/class/leds/beaglebone\:green\:usr0/trigger
		echo heartbeat > /sys/class/leds/beaglebone\:green\:usr1/trigger
		echo heartbeat > /sys/class/leds/beaglebone\:green\:usr2/trigger
		echo heartbeat > /sys/class/leds/beaglebone\:green\:usr3/trigger
	fi
	echo "-----------------------------"
	flush_cache
	umount $(dev2dir ${destination}p1) > /dev/null 2>&1 || true
	umount $(dev2dir ${destination}p2) > /dev/null 2>&1 || true
	inf_loop
}

check_eeprom () {

	eeprom="/sys/bus/i2c/devices/0-0050/eeprom"

	#Flash BeagleBone Black's eeprom:
	eeprom_location=$(ls /sys/devices/ocp.*/44e0b000.i2c/i2c-0/0-0050/eeprom 2> /dev/null)
	eeprom_header=$(hexdump -e '8/1 "%c"' ${eeprom} -s 5 -n 3)
	if [ "x${eeprom_header}" = "x335" ] ; then
		echo "Valid EEPROM header found"
	else
		echo "Invalid EEPROM header detected"
		if [ -f /opt/scripts/device/bone/bbb-eeprom.dump ] ; then
			if [ ! "x${eeprom_location}" = "x" ] ; then
				echo "Adding header to EEPROM"
				dd if=/opt/scripts/device/bone/bbb-eeprom.dump of=${eeprom_location}
				sync

				#We have to reboot, as the kernel only loads the eMMC cape
				# with a valid header
				reboot -f

				#We shouldnt hit this...
				exit
			fi
		fi
	fi
}

check_running_system () {
	if [ ! -f /boot/uboot/uEnv.txt ] ; then
		echo "Error: script halting, system unrecognized..."
		echo "unable to find: [/boot/uboot/uEnv.txt] is ${source}p1 mounted?"
		exit 1
	fi

	echo "-----------------------------"
	echo "debug copying: [${source}] -> [${destination}]"
	lsblk
	echo "-----------------------------"

	if [ ! -b "${destination}" ] ; then
		echo "Error: [${destination}] does not exist"
		write_failure
	fi

	#/ is ro...
	#if [ -L /etc/mtab && -r /proc/mounts ] ; then
	#	rm /etc/mtab && ln -s /proc/mounts /etc/mtab
	#fi
}

fdisk_toggle_boot () {
	fdisk ${destination} <<-__EOF__
	a
	1
	w
	__EOF__
	flush_cache
}

format_boot () {
	LC_ALL=C fdisk -l ${destination} | grep ${destination}p1 | grep '*' || fdisk_toggle_boot

	mkfs.vfat -F 16 ${destination}p1 -n BOOT
	flush_cache
}

format_root () {
	mkfs.ext4 ${destination}p2 -L rootfs
	flush_cache
}

partition_drive () {
	flush_cache
	dd if=/dev/zero of=${destination} bs=1M count=108
	sync
	dd if=${destination} of=/dev/null bs=1M count=108
	sync
	flush_cache

	#96Mb fat formatted boot partition
	LC_ALL=C sfdisk --force --in-order --Linux --unit M "${destination}" <<-__EOF__
		1,96,0xe,*
		,,,-
	__EOF__

	flush_cache
	format_boot
	format_root
}

copy_boot () {
	mkdir -p /tmp/boot/ || true
	mount ${destination}p1 /tmp/boot/ -o sync

	rsync -aAX /boot/uboot/ /tmp/boot/ || write_failure
	flush_cache

	umount /tmp/boot/ || umount -l /tmp/boot/ || write_failure
}

copy_rootfs () {
	mkdir -p /tmp/rootfs/ || true
	mount ${destination}p2 /tmp/rootfs/ -o async,noatime

	echo "rsync: / -> /tmp/rootfs/"
	rsync -aAX /* /tmp/rootfs/ --exclude={/dev/*,/proc/*,/sys/*,/tmp/*,/run/*,/mnt/*,/media/*,/lost+found,/boot/*,/lib/modules/*} || write_failure
	flush_cache

	mkdir -p /tmp/rootfs/boot/uboot/ || true
	mkdir -p /tmp/rootfs/lib/modules/$(uname -r)/ || true

	echo "rsync: /lib/modules/$(uname -r)/ -> /tmp/rootfs/lib/modules/$(uname -r)/"
	rsync -aAX /lib/modules/$(uname -r)/* /tmp/rootfs/lib/modules/$(uname -r)/ || write_failure
	flush_cache

	umount /tmp/rootfs/ || umount -l /tmp/rootfs/ || write_failure

	#https://github.com/beagleboard/meta-beagleboard/blob/master/contrib/bone-flash-tool/emmc.sh#L158-L159
	# force writeback of eMMC buffers
	dd if=${destination} of=/dev/null count=100000

	echo ""
	echo "This script has now completed it's task"
	echo "-----------------------------"

}

check_running_system
partition_drive
copy_boot
copy_rootfs
