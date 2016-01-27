from os.path import dirname, join

import convertible
from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.tools.hardware import Hardware
from syncloud_platform.tools.mount import Mount

DIR = dirname(__file__)


CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

logger.init(console=True)

def get_hardware():
    config = PlatformConfig(CONFIG_DIR)
    return Hardware(config, None, Mount(config))

def test_list():
    disks = get_hardware().available_disks(open(join(DIR, 'hardware', 'lsblk')).read())
    assert len(disks) == 3
    assert len(disks[0].partitions) == 1
    assert disks[1].partitions[2].mount_point == '/opt/disk/external'
    assert len(disks[1].partitions) == 3


def test_loop_support():
    disks = get_hardware().available_disks(open(join(DIR, 'hardware', 'lsblk')).read())
    assert len(disks) == 3
    assert len(disks[2].partitions) == 1
    convertible.to_dict(disks)
