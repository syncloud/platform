from os import unlink
import os
from os.path import islink, join
from subprocess import check_output
from os import path
from syncloud_app import logger
from syncloud_platform.systemd import systemctl
from syncloud_platform.tools.chown import chown


class Hardware:

    def __init__(self, platform_config, event_trigger, mount, lsblk, path_checker):
        self.platform_config = platform_config
        self.event_trigger = event_trigger
        self.mount = mount
        self.lsblk = lsblk
        self.path_checker = path_checker
        self.log = logger.get_logger('hardware')

    def available_disks(self):
        return self.lsblk.available_disks()

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
            self.log.info('external disk is mounted')
            permissions_support = external_mount.permissions_support()
        else:
            self.log.info('internal mount')
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
        self.log.info('checking external disk')
        if self.path_checker.external_disk_link_exists() and not self.lsblk.is_external_disk_attached():
            self.deactivate_disk()
