from os.path import dirname, join

from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.disks.lsblk import Lsblk
from syncloud_platform.disks.path_checker import PathChecker


DIR = dirname(__file__)

CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

logger.init(console=True)

default_output = '''NAME="/dev/sda" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06"
NAME="/dev/sda1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/abc" PARTTYPE="0x83" FSTYPE="ntfs" MODEL=""
NAME="/dev/sda2" SIZE="1K" TYPE="part" MOUNTPOINT="" PARTTYPE="0x5" FSTYPE="" MODEL=""
NAME="/dev/sda3" SIZE="5.0G" TYPE="part" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL=""
NAME="/dev/sda5" SIZE="7.4G" TYPE="part" MOUNTPOINT="[SWAP]" PARTTYPE="0x83" FSTYPE="" MODEL=""
NAME="/dev/sdb" SIZE="232.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="TOSHIBA MK2552GS"
NAME="/dev/sdb1" SIZE="100M" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL=""
NAME="/dev/sdb2" SIZE="48.7G" TYPE="part" MOUNTPOINT="" PARTTYPE="0x7" FSTYPE="" MODEL=""
NAME="/dev/sdb3" SIZE="184.1G" TYPE="part" MOUNTPOINT="/opt/disk/external" PARTTYPE="0x83" FSTYPE="ext4" MODEL=""
NAME="/dev/sdc" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06"
NAME="/dev/sdc1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ntfs" MODEL=""
NAME="/dev/sdc2" SIZE="1K" TYPE="part" MOUNTPOINT="" PARTTYPE="0x5" FSTYPE="" MODEL=""
NAME="/dev/sr0" SIZE="3.4G" TYPE="rom" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="CDDVDW SN-208AB "
NAME="/dev/sdc" SIZE="55.9G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="INTEL SSDSC2CW06"
NAME="/dev/sdc1" SIZE="48.5G" TYPE="part" MOUNTPOINT="/test" PARTTYPE="0x83" FSTYPE="vfat" MODEL=" "
NAME="/dev/loop0" SIZE="10M" TYPE="loop" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL=""
NAME="/dev/sda" SIZE="3.7G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="BLANK DISK"
NAME="/dev/loop1" SIZE="41.1M" TYPE="loop" MOUNTPOINT="/snap/platform/180821" PARTTYPE="" FSTYPE="squashfs" MODEL=""'''


def test_list():

    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(default_output)
    assert len(disks) == 5
    assert len(disks[0].partitions) == 4
    assert disks[1].partitions[2].mount_point == '/opt/disk/external'
    assert len(disks[1].partitions) == 3


def test_loop_support():
    output = 'NAME="/dev/loop0" SIZE="10M" TYPE="loop" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL=""'
    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(output)
    assert len(disks) == 1


def test_empty_disk():

    output = 'NAME="/dev/sda" SIZE="3.7G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL="BLANK DISK"'

    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(output)

    assert len(disks[0].partitions) == 0


def test_do_not_show_squashfs():

    output = 'NAME="/dev/loop1" SIZE="41.1M" TYPE="loop" MOUNTPOINT="/snap/platform/180821" PARTTYPE="" FSTYPE="squashfs" MODEL=""'

    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(output)

    assert len(disks) == 0


def test_do_not_show_internal_disks():
    lsblk_output = 'NAME="/dev/mmcblk0" SIZE="14.4G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL=""\n'
    lsblk_output += 'NAME="/dev/mmcblk0p2" SIZE="14.4G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ext4" MODEL=""\n'
    lsblk_output += 'NAME="/dev/mmcblk0p1" SIZE="41.8M" TYPE="part" MOUNTPOINT="" PARTTYPE="0xc" FSTYPE="vfat" MODEL=""'

    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(lsblk_output)

    assert len(disks) == 0


def test_do_not_show_disks_with_root_partition():
    lsblk_output = 'NAME="/dev/sdb" SIZE="14.4G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL=""\n'
    lsblk_output += 'NAME="/dev/sdb1" SIZE="14.4G" TYPE="part" MOUNTPOINT="/" PARTTYPE="0x83" FSTYPE="ext4" MODEL=""\n'
    lsblk_output += 'NAME="/dev/sdb2" SIZE="41.8M" TYPE="part" MOUNTPOINT="" PARTTYPE="0xc" FSTYPE="vfat" MODEL=""'

    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(lsblk_output)

    assert len(disks) == 0


def test_default_empty_disk_name():
    lsblk_output = 'NAME="/dev/sdb" SIZE="14.4G" TYPE="disk" MOUNTPOINT="" PARTTYPE="" FSTYPE="" MODEL=""\n'
    lsblk_output += 'NAME="/dev/sdb1" SIZE="14.4G" TYPE="part" MOUNTPOINT="" PARTTYPE="0x83" FSTYPE="ext4" MODEL=""\n'
    lsblk_output += 'NAME="/dev/sdb2" SIZE="41.8M" TYPE="part" MOUNTPOINT="" PARTTYPE="0xc" FSTYPE="vfat" MODEL=""'

    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    disks = lsblk.available_disks(lsblk_output)
    assert disks[0].name == 'Disk'


def test_is_external_disk_attached():
    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))
    assert lsblk.is_external_disk_attached(default_output, '/opt/disk/external')


def test_is_external_disk_detached():
    platform_config = PlatformConfig(CONFIG_DIR)
    lsblk = Lsblk(platform_config, PathChecker(platform_config))

    assert not lsblk.is_external_disk_attached(default_output, '/opt/disk/detached')
