from os import unlink
import os
from os.path import islink, join
import re
from subprocess import check_output
from os import path
from syncloud_app import logger
from syncloud_platform.systemd import systemctl
from syncloud_platform.tools.chown import chown

PARTTYPE_EXTENDED = '0x5'


class Hardware:

    def __init__(self, platform_config, event_trigger, mount):
        self.platform_config = platform_config
        self.event_trigger = event_trigger
        self.mount = mount
        self.log = logger.get_logger('hardware')

    def available_disks(self, lsblk_output=None):
        if not lsblk_output:
            lsblk_output = check_output('lsblk -Pp -o NAME,SIZE,TYPE,MOUNTPOINT,PARTTYPE,MODEL', shell=True)
        disks = []
        disk = None
        for line in lsblk_output.splitlines():
            match = re.match(
                r'NAME="(.*)" SIZE="(.*)" TYPE="(.*)" MOUNTPOINT="(.*)" PARTTYPE="(.*)" MODEL="(.*)"',
                line.strip())

            lsblk = LsblkEntry(match.group(1), match.group(2), match.group(3),
                               match.group(4), match.group(5), match.group(6).strip())

            if lsblk.type in ('disk', 'loop'):
                disk = Disk(lsblk.model.split(' ')[0])
                if lsblk.type == 'loop':
                    disk.add_partiotion(lsblk, self.platform_config)
                disks.append(disk)

            elif lsblk.type == 'part':
                disk.add_partiotion(lsblk, self.platform_config)

        disks_with_partitions = [d for d in disks if d.partitions]
        return disks_with_partitions

    def activate_disk(self, device):
        self.log.info('activate disk: {0}'.format(device))
        self.deactivate_disk()

        check_output('udisksctl mount -b {0}'.format(device), shell=True)
        mount_entry = self.mount.mounted_disk_by_device(device)
        check_output('udisksctl unmount -b {0}'.format(device), shell=True)
        systemctl.add_mount(mount_entry)

        self.relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_external_disk_dir())

    def deactivate_disk(self):
        self.log.info('deactivate disk')
        self.relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_internal_disk_dir())
        systemctl.remove_mount()

    def init_app_storage(self, app_id, owner=None):
        external_mount = self.mount.get_mounted_external_disk()
        if external_mount:
            permissions_support = external_mount.permissions_support()
        else:
            permissions_support = True

        app_storage_dir = join(self.platform_config.get_disk_link(), app_id)
        if not path.exists(app_storage_dir):
            os.mkdir(app_storage_dir)
        if owner and permissions_support:
            self.log.info('fixing permissions on {0}'.format(app_storage_dir))
            chown(owner, app_storage_dir)
        else:
            self.log.info('not fixing permissions')
        return app_storage_dir

    def relink_disk(self, link, target):

        os.chmod(target, 0755)

        if islink(link):
            unlink(link)
        os.symlink(target, link)

        self.event_trigger.trigger_app_event_disk(self.platform_config.apps_root())

    def check_external_disk(self):
        self.log.info('check external disk')
        if self.mount.external_disk_link_exists() and not self.mount.get_mounted_external_disk(): 
            self.deactivate_disk() 


class LsblkEntry:
    def __init__(self, name, size, type, mountpoint, parttype, model):
        self.name = name
        self.size = size
        self.type = type
        self.mountpoint = mountpoint
        self.parttype = parttype
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

    def add_partiotion(self, lsblk, platform_config):
        mountable = False
        mount_point = lsblk.mountpoint
        if not lsblk.is_extended_partition():
            if not mount_point or mount_point == platform_config.get_external_disk_dir():
                mountable = True

        if lsblk.is_boot_disk():
            mountable = False
        active = False
        if mount_point == platform_config.get_external_disk_dir() and self.external_disk_link_exists(platform_config):
            active = True
        if mountable:
            self.partitions.append(Partition(lsblk.size, lsblk.name, mount_point, active))

    def external_disk_link_exists(self, platform_config):
        return path.realpath(platform_config.get_disk_link()) == platform_config.get_external_disk_dir()

    def __str__(self):
        return '{0}: {1}'.format(self.name, ','.join(map(str, self.partitions)))

