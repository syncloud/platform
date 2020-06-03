from subprocess import check_output
import re
from syncloudlib import logger

PARTTYPE_EXTENDED = '0x5'


class Lsblk:

    def __init__(self, platform_config, path_checker):
        self.platform_config = platform_config
        self.path_checker = path_checker
        self.log = logger.get_logger('lsblk')

    def available_disks(self, lsblk_output=None):
        return [d for d in self.all_disks(lsblk_output) if not d.is_internal() and not d.has_root_partition()]

    def all_disks(self, lsblk_output=None):
        if not lsblk_output:
            lsblk_output = check_output('lsblk -Pp -o NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL', shell=True).decode()
        # self.log.info(lsblk_output)

        disks = {}
        for line in lsblk_output.splitlines():
            self.log.info('parsing line: {0}'.format(line))
            match = re.match(
                r'NAME="(.*)" SIZE="(.*)" TYPE="(.*)" MOUNTPOINT="(.*)" PARTTYPE="(.*)" FSTYPE="(.*)" MODEL="(.*)"',
                line.strip())

            lsblk_entry = LsblkEntry(match.group(1), match.group(2), match.group(3),
                                     match.group(4), match.group(5), match.group(6), match.group(7).strip())
            
            if lsblk_entry.is_supported_type() and lsblk_entry.is_supported_fs_type():
                device = lsblk_entry.name
                disk_name = lsblk_entry.model
                self.log.info('adding disk: {0}'.format(disk_name))
                disk = Disk(disk_name, device, lsblk_entry.size, [])
                if lsblk_entry.is_single_partition_disk():
                    self.log.info('adding single partition disk: {0}'.format(device))

                    disk.name = lsblk_entry.type
                    partition = self.create_partition(lsblk_entry)
                    disk.add_partition(partition)

                disks[device] = disk

            elif lsblk_entry.type == 'part':
                self.log.info('adding regular parrition: {0}'.format(lsblk_entry.name))
                partition = self.create_partition(lsblk_entry)
                parent_device = lsblk_entry.parent_device()
                disk = disks[parent_device]
                disk.add_partition(partition)

        return disks.values()

    def create_partition(self, lsblk_entry):
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

        return Partition(lsblk_entry.size, lsblk_entry.name, mount_point, active, lsblk_entry.fs_type, mountable)

    def is_external_disk_attached(self, lsblk_output=None, disk_dir=None):
        if not disk_dir:
            disk_dir = self.platform_config.get_external_disk_dir()
        for disk in self.all_disks(lsblk_output):
            for partition in disk.partitions:
                if partition.mount_point == disk_dir:
                    self.log.info('external disk is attached')
                    return True
        self.log.info('external disk is detached')
        return False

    def find_partition_by_device(self, device, lsblk_output=None):
        for disk in self.all_disks(lsblk_output):
            for partition in disk.partitions:
                if partition.device == device:
                    self.log.info('partition found')
                    return partition
        self.log.info('partition not found')
        return None

    def find_partition_by_dir(self, mount_dir, lsblk_output=None):
        for disk in self.all_disks(lsblk_output):
            for partition in disk.partitions:
                if partition.mount_point == mount_dir:
                    self.log.info('partition found')
                    return partition
        self.log.info('partition not found')
        return None


class LsblkEntry:
    def __init__(self, name, size, device_type, mount_point, part_type, fs_type, model):
        self.name = name
        self.size = size
        self.type = device_type
        self.mount_point = mount_point
        self.parttype = part_type
        self.fs_type = fs_type
        self.model = model

    def is_extended_partition(self):
        return self.parttype == PARTTYPE_EXTENDED

    def is_boot_disk(self):
        return '/dev/mmcblk0' in self.name

    def is_supported_type(self):
        if self.type in ('disk', 'loop'):
            return True
        if self.type.startswith('raid'):
            return True
        return False

    def is_supported_fs_type(self):
        if self.fs_type == 'squashfs':
            return False
        if self.fs_type == 'linux_raid_member':
            return False
        return True

    def is_single_partition_disk(self):
        if self.type == 'loop':
            return True
        if self.type.startswith('raid'):
            return True
        return False

    def parent_device(self):
        match = re.match(r'(.*?)p?\d*$', self.name)
        return match.group(1)


class Partition:
    def __init__(self, size, device, mount_point, active, fs_type, mountable):
        self.size = size
        self.device = device
        self.mount_point = mount_point
        self.active = active
        self.fs_type = fs_type
        self.mountable = mountable
        self.extendable = False

    def permissions_support(self):
        return self.fs_type not in ['vfat', 'exfat']

    def is_root_fs(self):
        return self.mount_point == '/'

    def __str__(self):
        return '{0}, {1}, {2}, {3}'.format(self.device, self.size, self.mount_point, self.active)


class Disk:
    def __init__(self, name, device, size, partitions):
        
        if name == '':
            name = 'Disk'
        self.name = name
        self.partitions = partitions
        self.device = device
        self.size = size
        self.active = False
    
    def is_internal(self):
        return self.device.startswith('/dev/mmcblk')

    def has_root_partition(self):
        return self.find_root_partition() is not None

    def add_partition(self, partition):
        if partition.active:
            self.active = True
        self.partitions.append(partition)

    def find_root_partition(self):
        return next((p for p in self.partitions if p.is_root_fs()), None)

    def __str__(self):
        return '{0}: {1}'.format(self.name, ','.join(list(map(str, self.partitions))))
