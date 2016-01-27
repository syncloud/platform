from subprocess import check_output
from os import path
from syncloud_app import logger

class Mount:
    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.log = logger.get_logger('mount')

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
                    .replace('default_permissions', 'permissions')\
                    .replace('nodev,', '')
                return MountEntry(device, dir, type, options)
        return None

    def get_mounted_external_disk(self):
        mount_point = None

        if self.external_disk_link_exists():
            disk_dir = self.platform_config.get_external_disk_dir()
            mount_entry = self.mounted_disk_by_dir(disk_dir)

        return mount_point

    def external_disk_link_exists(self):
        return path.realpath(self.platform_config.get_disk_link()) == self.platform_config.get_external_disk_dir()


class MountEntry:

    def __init__(self, device, dir, type, options):
        self.device = device
        self.dir = dir
        self.type = type
        self.options = options
        self.log = logger.get_logger('mount_entry')
        self.log.info('entry: {0}, {1}, {2}, {3}'.format(device, dir, type, options))

    def permissions_support(self):
        return 'fat' not in self.type
