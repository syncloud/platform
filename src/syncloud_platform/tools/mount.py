from subprocess import check_output
from os import path
from syncloud_app import logger
import re

class Mount:
    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.log = logger.get_logger('mount')

    def mounted_disk_by_device(self, device, mount_output=None):
        self.log.info('searching by device: {0}'.format(device))
        return self.__mounted_disk(lambda entry: entry.startswith('{0} on'.format(device)), mount_output)

    def mounted_disk_by_dir(self, dir, mount_output=None):
        self.log.info('searching by dir: {0}'.format(dir))
        return self.__mounted_disk(lambda entry: ' on {0} type'.format(dir) in entry, mount_output)

    def __mounted_disk(self, entry_filter, mount_output=None):
        if not mount_output:
            mount_output = check_output('mount', shell=True)
        self.log.info('searching mount')
        for entry in mount_output.splitlines():
            self.log.info('entry: {0}'.format(entry))
            if entry_filter(entry):
                self.log.info('found')
                parts_on = entry.split(' on ')
                device = parts_on[0]
                parts_type = parts_on[1].split(' type ')
                dir = parts_type[0]
                parts_options = parts_type[1].split(' ')
                type = parts_options[0].replace('fuseblk', 'ntfs')
                options = self.get_options(type, parts_options[1])
                return MountEntry(device, dir, type, options)
        self.log.info('not found')
        return None

    def get_options(self, type, udisks_options):
        options = udisks_options.strip('()')\
                    .replace('codepage=cp', 'codepage=')\
                    .replace('default_permissions', 'permissions')\
                    .replace('nodev,', '')
        options = re.sub('fmask=\d+', 'fmask=0000', options)
        options = re.sub('dmask=\d+', 'dmask=0000', options)
        options = re.sub('umask=\d+', 'umask=0000', options)
        return options

    def get_mounted_external_disk(self):
        mount_point = None

        if self.external_disk_link_exists():
            disk_dir = self.platform_config.get_external_disk_dir()
            mount_point = self.mounted_disk_by_dir(disk_dir)

        return mount_point

    def external_disk_link_exists(self):
        real_link_path = path.realpath(self.platform_config.get_disk_link())
        self.log.info('real link path: {0}'.format(real_link_path))
        external_disk_path = self.platform_config.get_external_disk_dir()
        self.log.info('external disk path: {0}'.format(external_disk_path))
        return real_link_path == external_disk_path


class MountEntry:

    def __init__(self, device, dir, type, options):
        self.device = device
        self.dir = dir
        self.type = type
        self.options = options
        self.log = logger.get_logger('mount_entry')
        self.log.info('entry: {0}, {1}, {2}, {3}'.format(device, dir, type, options))

    def permissions_support(self):
        supported = 'fat' not in self.type
        self.log.info('permissions support: {0}'.format(supported))
        return supported
