package storage

import "github.com/syncloud/platform/config"

type Lsblk struct {

}

func NewLsblk(systemConfig config.SystemConfig, path_checker) {
	self.platform_config = platform_config
	self.path_checker = path_checker
	self.log = logger.get_logger('lsblk')
}

def available_disks(self, lsblk_output=None):
return [d for d in self.all_disks(lsblk_output) if not d.is_internal() and not d.has_root_partition()]

def all_disks(self, lsblk_output=None):
if not lsblk_output:
lsblk_output = check_output('lsblk -Pp -o NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,FSTYPE,MODEL', shell=True).decode()
self.log.info(lsblk_output)

disks = {}
for line in lsblk_output.splitlines():
if not line.strip():
continue

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
self.log.info('adding regular partition: {0}'.format(lsblk_entry.name))
partition = self.create_partition(lsblk_entry)
parent_device = lsblk_entry.parent_device()
if parent_device in disks:
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

return Partition(lsblk_entry.size, lsblk_entry.name, mount_point, active, lsblk_entry.get_fs_type(), mountable)

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


