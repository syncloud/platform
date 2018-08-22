from os.path import dirname, join

import convertible
from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.disks.lsblk import Lsblk
from syncloud_platform.disks.path_checker import PathChecker


DIR = dirname(__file__)
LSBLK = open(join(DIR, 'hardware', 'lsblk')).read()


CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

logger.init(console=True)


def get_lsblk():
    platform_config = PlatformConfig(CONFIG_DIR)
    return Lsblk(platform_config, PathChecker(platform_config))


def test_list():
    disks = get_lsblk().available_disks(LSBLK)
    assert len(disks) == 5
    assert len(disks[0].partitions) == 4
    assert disks[1].partitions[2].mount_point == '/opt/disk/external'
    assert len(disks[1].partitions) == 3


def test_loop_support():
    disks = get_lsblk().available_disks(LSBLK)
    assert len(disks) == 5
    assert len(disks[2].partitions) == 1
    convertible.to_dict(disks)


def test_empty_disk():
    disks = get_lsblk().available_disks(LSBLK)
    assert len(disks[4].partitions) == 0


def test_is_external_disk_attached():
    assert get_lsblk().is_external_disk_attached(open(join(DIR, 'hardware', 'lsblk')).read(), '/opt/disk/external')


def test_is_external_disk_detached():
    assert not get_lsblk().is_external_disk_attached(open(join(DIR, 'hardware', 'lsblk')).read(), '/opt/disk/detached')
