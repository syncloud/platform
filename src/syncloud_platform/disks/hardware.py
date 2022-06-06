import os
from os import path
from os import unlink
from os.path import islink, join, isdir

from syncloudlib import logger, fs

supported_fs = {
    'ext2',
    'ext3',
    'ext4',
    'raid'
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
