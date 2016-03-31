from subprocess import check_output
from syncloud_app import logger

known_fs_options = {
    'ext4': 'rw,nosuid,relatime,data=ordered,uhelper=udisks2',
    'vfat': 'rw,nosuid,relatime,fmask=0000,dmask=0000,codepage=437,iocharset=iso8859-1,shortname=mixed,showexec,utf8,flush,errors=remount-ro',
    'ntfs': 'rw,nosuid,relatime,user_id=0,group_id=0,permissions,allow_other,blksize=4096',
    'exfat': 'rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other,blksize=4096',
    'ext2': 'rw,nosuid,nodev,realtime'
}


class Mount:
    def __init__(self, platform_config, path_checker, lsblk):
        self.platform_config = platform_config
        self.path_checker = path_checker
        self.lsblk = lsblk
        self.log = logger.get_logger('mount')

    def mounted_disk_by_device(self, device, mount_output=None, lsblk_output=None):
        self.log.info('searching by device: {0}'.format(device))
        partition = self.lsblk.find_partition_by_device(device, lsblk_output)
        if not partition:
            return
        return self.__mounted_disk(lambda entry: entry.startswith('{0} on'.format(device)), partition.fs_type, mount_output)

    def mounted_disk_by_dir(self, dir, mount_output=None, lsblk_output=None):
        self.log.info('searching by dir: {0}'.format(dir))
        partition = self.lsblk.find_partition_by_dir(dir, lsblk_output)
        if not partition:
            return
        return self.__mounted_disk(lambda entry: ' on {0} type'.format(dir) in entry, partition.fs_type, mount_output)

    def __mounted_disk(self, entry_filter, fs_type, mount_output=None):

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
                options = self.get_options(fs_type, parts_options[1])
                return MountEntry(device, dir, fs_type, options)
        self.log.info('not found')
        return None

    def get_options(self, type, udisks_options):

        if type in known_fs_options:
            return known_fs_options[type]

        return udisks_options.strip('()')

    def get_mounted_external_disk(self):
        mount_point = None

        if self.path_checker.external_disk_link_exists():
            disk_dir = self.platform_config.get_external_disk_dir()
            mount_point = self.mounted_disk_by_dir(disk_dir)

        return mount_point


class MountEntry:

    def __init__(self, device, dir, type, options):
        self.device = device
        self.dir = dir
        self.type = type
        self.options = options
        self.log = logger.get_logger('mount_entry')
        self.log.info('entry: {0}, {1}, {2}, {3}'.format(device, dir, type, options))

    def permissions_support(self):
        supported = self.type not in ['vfat', 'exfat']
        self.log.info('permissions support: {0}'.format(supported))
        return supported
