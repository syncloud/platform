from os.path import dirname, join

import convertible
from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.tools.hardware import Hardware

DIR = dirname(__file__)


CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

logger.init(console=True)


def test_list():
    disks = Hardware(PlatformConfig(CONFIG_DIR), None).available_disks(open(join(DIR, 'hardware', 'lsblk')).read())
    assert len(disks) == 3
    assert len(disks[0].partitions) == 1
    assert disks[1].partitions[2].mount_point == '/opt/disk/external'
    assert len(disks[1].partitions) == 3


def test_loop_support():
    disks = Hardware(PlatformConfig(CONFIG_DIR), None).available_disks(open(join(DIR, 'hardware', 'lsblk')).read())
    assert len(disks) == 3
    assert len(disks[2].partitions) == 1
    convertible.to_dict(disks)


def test_get_mount_info_by_device():
    mount_point = Hardware(PlatformConfig(CONFIG_DIR), None).mounted_disk_by_device('/dev/sdc1', open(join(DIR, 'hardware', 'mount')).read())
    assert mount_point.device == '/dev/sdc1'
    assert mount_point.dir == '/media/root/long name'
    assert mount_point.type == 'vfat'
    assert mount_point.options == 'rw,nosuid,nodev,relatime,fmask=0022,dmask=0077,codepage=437,iocharset=iso8859-1,' \
                                  'shortname=mixed,showexec,utf8,flush,errors=remount-ro'


def test_get_mount_info_by_dir():
    mount_point = Hardware(PlatformConfig(CONFIG_DIR), None).mounted_disk_by_dir('/opt/disk/external', open(join(DIR, 'hardware', 'mount')).read())
    assert mount_point.device == '/dev/sdc2'
    assert mount_point.dir == '/opt/disk/external'
    assert mount_point.type == 'ext4'
    assert mount_point.options == 'rw,nosuid,nodev,relatime,data=ordered,uhelper=udisks2'


def test_ntfs_permissions():
    mount_point = Hardware(PlatformConfig(CONFIG_DIR), None).mounted_disk_by_device('/dev/sdc3', open(join(DIR, 'hardware', 'mount')).read())
    assert mount_point.device == '/dev/sdc3'
    assert mount_point.dir == '/media/ntfs'
    assert mount_point.type == 'ntfs'
    assert mount_point.options == 'rw,nosuid,nodev,relatime,user_id=0,group_id=0,permissions,allow_other,' \
                                  'blksize=4096'
