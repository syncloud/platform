from syncloud_app import logger

from syncloud_platform.disks.lsblk import Disk, Partition

logger.init(console=True)


def test_find_root_partition_some():

    disk = Disk('disk', '/dev/sda', 20, [
        Partition(10, '/dev/sda1', '/', True, 'ext4', False),
        Partition(10, '/dev/sda2', '', True, 'ext4', True)
    ])

    assert disk.find_root_partition().device == '/dev/sda1'


def test_find_root_partition_none():

    disk = Disk('disk', '/dev/sda', 20, [
        Partition(10, '/dev/sda1', '/my', True, 'ext4', False),
        Partition(10, '/dev/sda2', '', True, 'ext4', True)
    ])

    assert disk.find_root_partition() is None
