from os import unlink
import os
from os.path import islink, join, isdir
from os import path
from syncloudlib import logger, fs
from syncloud_platform.disks.lsblk import Partition
from syncloud_platform.gaplib.linux import parted

from syncloud_platform.rest.service_exception import ServiceException

supported_fs_options = {
    # 'vfat': 'rw,nosuid,relatime,fmask=0000,dmask=0000,codepage=437,iocharset=iso8859-1,'
    #         'shortname=mixed,showexec,utf8,flush,errors=remount-ro',
    # 'ntfs': 'rw,nosuid,relatime,user_id=0,group_id=0,permissions,allow_other,blksize=4096',
    # 'exfat': 'rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other,blksize=4096',
    'ext2': 'rw,nosuid,nodev,relatime',
    'ext3': 'rw,nosuid,nodev,relatime',
    'ext4': 'rw,nosuid,relatime,data=ordered'
}

EXTENDABLE_FREE_PERCENT = 10


def has_unallocated_space_at_the_end(parted_output):
    last_line = parted_output.splitlines()[-1]
    if 'free' not in last_line:
        return False
    free = float(last_line.split(':')[3][:-1])
    return free > EXTENDABLE_FREE_PERCENT


class Hardware:

    def __init__(self, platform_config, event_trigger, lsblk, path_checker, systemctl):
        self.platform_config = platform_config
        self.systemctl = systemctl
        self.event_trigger = event_trigger
        self.lsblk = lsblk
        self.path_checker = path_checker
        self.log = logger.get_logger('hardware')

    def available_disks(self):
        return [d for d in self.lsblk.available_disks() if not d.is_internal()]

    def root_partition(self):
        disks = self.lsblk.all_disks()

        partition = Partition(0, 'unknown', '/', True, 'unknown', False)

        boot_disk = next((d for d in disks if d.find_root_partition() is not None), None)
        if boot_disk:
            partition = boot_disk.find_root_partition()
            parted_output = parted(boot_disk.device)
            partition.extendable = has_unallocated_space_at_the_end(parted_output)

        return partition

    def activate_disk(self, device):
        self.log.info('activate disk: {0}'.format(device))
        self.deactivate_disk()

        partition = self.lsblk.find_partition_by_device(device)
        if not partition:
            error_message = 'unable to find device: {0}'.format(device)
            self.log.error(error_message)
            raise Exception(error_message)

        fs_type = partition.fs_type
        if fs_type not in supported_fs_options:
            error_message = 'Filesystem type is not supported: {0}' \
                            ', use on of the following: {1}'.format(fs_type, supported_fs_options.keys())
            self.log.error(error_message)
            raise ServiceException(error_message)

        self.systemctl.add_mount(device, fs_type, supported_fs_options[fs_type])

        relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_external_disk_dir())
        self.event_trigger.trigger_app_event_disk()

    def deactivate_disk(self):
        self.log.info('deactivate disk')
        relink_disk(
            self.platform_config.get_disk_link(),
            self.platform_config.get_internal_disk_dir())
        self.event_trigger.trigger_app_event_disk()    
        self.systemctl.remove_mount()

    def get_app_storage_dir(self, app_id):
        app_storage_dir = join(self.platform_config.get_disk_link(), app_id)
        return app_storage_dir

    def init_app_storage(self, app_id, owner=None):
        app_storage_dir = self.get_app_storage_dir(app_id)
        if not path.exists(app_storage_dir):
            os.mkdir(app_storage_dir)
        if owner:
            self.log.info('fixing permissions on {0}'.format(app_storage_dir))
            fs.chownpath(app_storage_dir, owner, recursive=True)
        else:
            self.log.info('not fixing permissions')
        return app_storage_dir

    def init_disk(self):

        if not isdir(self.platform_config.get_disk_root()):
            os.mkdir(self.platform_config.get_disk_root())

        if not isdir(self.platform_config.get_internal_disk_dir()):
            os.mkdir(self.platform_config.get_internal_disk_dir())

        if not self.path_checker.external_disk_link_exists():
            relink_disk(self.platform_config.get_disk_link(), self.platform_config.get_internal_disk_dir())


def relink_disk(link, target):

    os.chmod(target, 0o755)

    if islink(link):
        unlink(link)
    os.symlink(target, link)
