from os import unlink
import os
from os.path import islink, join
import re
from subprocess import check_output
from os import path
from syncloud_app import logger
from syncloud_platform.config.config import PLATFORM_CONFIG_DIR, PlatformConfig
from syncloud_platform.systemd import systemctl
from syncloud_platform.tools.chown import chown
from syncloud_platform.tools.events import trigger_app_event_disk

PARTTYPE_EXTENDED = '0x5'


class Hardware:

    def __init__(self, config_path=PLATFORM_CONFIG_DIR):
        self.platform_config = PlatformConfig(config_path)
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

    def mounted_disk_by_device(self, device, mount_output=None):
        return self.__mounted_disk(lambda entry: entry.startswith('{0} on'.format(device)), mount_output)

    def mounted_disk_by_dir(self, dir, mount_output=None):
        return self.__mounted_disk(lambda entry: ' on {0} type'.format(dir) in entry, mount_output)

    def __mounted_disk(self, entry_filter, mount_output=None):
        if not mount_output:
            mount_output = check_output('mount', shell=True)
        for entry in mount_output.splitlines():
            if entry_filter(entry):
                parts_on = entry.split(' on ')
                device = parts_on[0]
                parts_type = parts_on[1].split(' type ')
                dir = parts_type[0]
                parts_options = parts_type[1].split(' ')
                type = parts_options[0].replace('fuseblk', 'ntfs')
                options = parts_options[1].strip('()')\
                    .replace('codepage=cp', 'codepage=')\
                    .replace('default_permissions', 'permissions')
                return MountEntry(device, dir, type, options)
        return None

    def activate_disk(self, device):

        self.deactivate_disk()

        check_output('udisksctl mount -b {0}'.format(device), shell=True)
        mount_entry = self.mounted_disk_by_device(device)
        check_output('udisksctl unmount -b {0}'.format(device), shell=True)
        systemctl.add_mount(mount_entry)

        self.relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_external_disk_dir())

    def deactivate_disk(self):
        self.relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_internal_disk_dir())
        systemctl.remove_mount()

    def init_app_storage(self, app_id, owner=None):

        if self.external_disk_is_mounted():
            path.realpath(self.platform_config.get_disk_link())
            disk_dir = self.platform_config.get_external_disk_dir()
            mount_entry = self.mounted_disk_by_dir(disk_dir)
            permissions_support = mount_entry.permissions_support()
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

        if islink(link):
            unlink(link)
        os.symlink(target, link)

        trigger_app_event_disk(self.platform_config.apps_root())

    def external_disk_is_mounted(self):
        return path.realpath(self.platform_config.get_disk_link()) == self.platform_config.get_external_disk_dir()


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
        if mount_point == platform_config.get_external_disk_dir() and self.external_disk_is_mounted(platform_config):
            active = True
        if mountable:
            self.partitions.append(Partition(lsblk.size, lsblk.name, mount_point, active))

    def external_disk_is_mounted(self, platform_config):
        return path.realpath(platform_config.get_disk_link()) == platform_config.get_external_disk_dir()


class MountEntry:

    def __init__(self, device, dir, type, options):
        self.device = device
        self.dir = dir
        self.type = type
        self.options = options

    def permissions_support(self):
        return 'fat' not in self.type
