from subprocess import check_output
import re
from syncloud_app import logger

PARTTYPE_EXTENDED = '0x5'


class Lsblk:

    def __init__(self, platform_config, path_checker):
        self.platform_config = platform_config
        self.path_checker = path_checker
        self.log = logger.get_logger('lsblk')

    def available_disks(self, lsblk_output=None):
        if not lsblk_output:
            lsblk_output = check_output('lsblk -Pp -o NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,MODEL', shell=True)
        disks = []
        disk = None
        for line in lsblk_output.splitlines():
            match = re.match(
                r'NAME="(.*)" SIZE="(.*)" TYPE="(.*)" MOUNTPOINT="(.*)" PARTTYPE="(.*)" MODEL="(.*)"',
                line.strip())

            lsblk_entry = LsblkEntry(match.group(1), match.group(2), match.group(3),
                                     match.group(4), match.group(5), match.group(6).strip())

            if lsblk_entry.type in ('disk', 'loop'):
                disk_name = lsblk_entry.model.split(' ')[0]
                disk = Disk(disk_name)
                if lsblk_entry.type == 'loop':
                    self.add_partition(disk, lsblk_entry)
                disks.append(disk)

            elif lsblk_entry.type == 'part':
                self.add_partition(disk, lsblk_entry)

        disks_with_partitions = [d for d in disks if d.partitions]
        return disks_with_partitions

    def add_partition(self, disk, lsblk_entry):
        mountable = False
        mount_point = lsblk_entry.mount_point
        if not lsblk_entry.is_extended_partition():
            if not mount_point or mount_point == self.platform_config.get_external_disk_dir():
                mountable = True

        if lsblk_entry.is_boot_disk():
            mountable = False
        active = False
        if mount_point == self.platform_config.get_external_disk_dir() \
                and self.path_checker.external_disk_link_exists():
            active = True
        if mountable:
            disk.partitions.append(Partition(lsblk_entry.size, lsblk_entry.name, mount_point, active))

    def is_external_disk_attached(self, lsblk_output=None, disk_dir=None):
        if not disk_dir:
            disk_dir = self.platform_config.get_external_disk_dir()
        for disk in self.available_disks(lsblk_output):
            for partition in disk.partitions:
                if partition.mount_point == disk_dir:
                    self.log.info('external disk is attached')
                    return True
        self.log.info('external disk is detached')
        return False


class LsblkEntry:
    def __init__(self, name, size, fs_type, mount_point, part_type, model):
        self.name = name
        self.size = size
        self.type = fs_type
        self.mount_point = mount_point
        self.parttype = part_type
        self.model = model

    def is_extended_partition(self):
        return self.parttype == PARTTYPE_EXTENDED

    def is_boot_disk(self):
        return '/dev/mmcblk0' in self.name


class Partition:
    def __init__(self, size, device, mount_point, active):
        self.size = size
        self.device = device
        self.mount_point = mount_point
        self.active = active

    def __str__(self):
        return '{0}, {1}, {2}, {3}'.format(self.device, self.size, self.mount_point, self.active)


class Disk:
    def __init__(self, name):
        self.partitions = []
        self.name = name

    def __str__(self):
        return '{0}: {1}'.format(self.name, ','.join(map(str, self.partitions)))
