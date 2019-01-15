from os import path
from syncloudlib import logger


class PathChecker:

    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.log = logger.get_logger('path_checker')

    def external_disk_link_exists(self):
        real_link_path = path.realpath(self.platform_config.get_disk_link())
        self.log.info('real link path: {0}'.format(real_link_path))

        external_disk_path = self.platform_config.get_external_disk_dir()
        self.log.info('external disk path: {0}'.format(external_disk_path))

        link_exists = real_link_path == external_disk_path
        self.log.info('link exists: {0}'.format(link_exists))

        return link_exists
