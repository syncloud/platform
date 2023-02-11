import os
from os import unlink
from os.path import islink, isdir


class Hardware:

    def __init__(self, platform_config, path_checker):
        self.platform_config = platform_config
        self.path_checker = path_checker

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
