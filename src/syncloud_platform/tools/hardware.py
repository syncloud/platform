import json
from os import unlink
import os
from os.path import islink, join
import re
from subprocess import check_output
from syncloud_app import logger
from syncloud_platform.config.config import PLATFORM_CONFIG_DIR, PlatformConfig
from syncloud_platform.systemd import systemctl
from syncloud_platform.tools.chown import chown
from syncloud_platform.tools.touch import touch

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
            fields = dict()
            for field in re.split(r'" ', line.strip()):
                key, value = field.split('=', 1)
                value = value[1:]
                fields[key] = value

            if fields['TYPE'] == 'disk':
                disk = Disk(fields['MODEL'].split(' ')[0])
                disks.append(disk)

            elif fields['TYPE'] == 'part':

                mountable = False
                mount_point = fields['MOUNTPOINT']
                if not fields['PARTTYPE'] == PARTTYPE_EXTENDED:
                    if not mount_point or mount_point == self.platform_config.get_external_disk_dir():
                        mountable = True

                if '/dev/mmcblk0' in fields['NAME']:
                    mountable = False

                if mountable:
                    disk.partitions.append(Partition(fields['SIZE'], fields['NAME'], mount_point))
        disks_with_partitions = [d for d in disks if d.partitions]
        return disks_with_partitions

    def mounted_disk(self, device, mount_output=None):
        if not mount_output:
            mount_output = check_output('mount', shell=True)
        for entry in mount_output.splitlines():
            if entry.startswith('{0} on'.format(device)):
                parts_on = entry.split(' on ')
                device = parts_on[0]
                parts_type = parts_on[1].split(' type ')
                dir = parts_type[0]
                parts_options = parts_type[1].split(' ')
                type = parts_options[0]
                return MountEntry(device, dir, type, parts_options[1].strip('()'))
        return None

    def activate_disk(self, device, fix_permissions=True):

        self.deactivate_disk()

        check_output('udisksctl mount -b {0}'.format(device), shell=True)
        mount_entry = self.mounted_disk(device)
        check_output('udisksctl unmount -b {0}'.format(device), shell=True)
        systemctl.add_mount(mount_entry)

        relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_external_disk_dir(),
            fix_permissions)

    def deactivate_disk(self):
        relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_internal_disk_dir())
        systemctl.remove_mount()


def relink_disk(link, target, fix_permissions=True):

    log = logger.get_logger('hardware.relink_disk')

    if islink(link):
        unlink(link)
    os.symlink(target, link)
    if fix_permissions:
        log.info('fixing permissions')
        # TODO: We need to come up with some generic way of giving access to different apps
        chown('owncloud', link)
    else:
        log.info('not fixing permissions')

    touch(join(link, '.ocdata'))


class Partition:
    def __init__(self, size, device, mount_point):
        self.size = size
        self.device = device
        self.mount_point = mount_point


class Disk:
    def __init__(self, name):
        self.partitions = []
        self.name = name


class MountEntry:

    def __init__(self, device, dir, type, options):
        self.device = device
        self.dir = dir
        self.type = type
        self.options = options
