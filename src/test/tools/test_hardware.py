from os.path import dirname, join
from syncloud_app import logger

from syncloud_platform.tools.hardware import Hardware

DIR = dirname(__file__)


CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

logger.init(console=True)


def test_list():
    disks = Hardware(CONFIG_DIR).available_disks(open(join(DIR, 'hardware', 'lsblk')).read())
    assert len(disks) == 2
    assert len(disks[0].partitions) == 1
    assert disks[1].partitions[2].mount_point == '/opt/disk/external'
    assert len(disks[1].partitions) == 3


def test_get_mount_info():
    mount_point = Hardware(CONFIG_DIR).mounted_disk('/dev/sdc1', open(join(DIR, 'hardware', 'mount')).read())
    assert mount_point.device == '/dev/sdc1'
    assert mount_point.dir == '/media/root/long name'
    assert mount_point.type == 'vfat'
    assert mount_point.options == 'rw,nosuid,nodev,relatime,fmask=0022,dmask=0077,codepage=437,iocharset=iso8859-1,' \
                                  'shortname=mixed,showexec,utf8,flush,errors=remount-ro'
