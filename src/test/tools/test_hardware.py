from os.path import dirname, join

from syncloud_platform.tools.hardware import Hardware

DIR = dirname(__file__)


def test_list():
    disks = Hardware().available_disks(
        open(join(DIR, 'hardware', 'lshw.json')).read(),
        open(join(DIR, 'hardware', 'mount')).read()
    )
    assert len(disks) == 3
    assert disks[0].partitions[0].mount_point == '/opt/disk'


def test_get_mount_info():
    mount_point = Hardware().mounted_disk('/dev/sdc1', open(join(DIR, 'hardware', 'mount')).read())
    assert mount_point.device == '/dev/sdc1'
    assert mount_point.dir == '/media/root/long name'
    assert mount_point.type == 'vfat'
    assert mount_point.options == 'rw,nosuid,nodev,relatime,fmask=0022,dmask=0077,codepage=437,iocharset=iso8859-1,' \
                                  'shortname=mixed,showexec,utf8,flush,errors=remount-ro'
