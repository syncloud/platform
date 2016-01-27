
import shutil
import os
from os.path import join
from syncloud_app import logger


udev_file_name = '90-syncloud.rules'
udev_dir = '/etc/udev/rules.d'


class Udev:
    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.log = logger.get_logger('udev')
        self.from_path = join(self.platform_config.config_dir(), 'udev',  udev_file_name)
        self.to_path = join(udev_dir, udev_file_name)

    def add(self):
        self.log.info('adding')
        shutil.copy(self.from_path, self.to_path)

    def remove(self):
        self.log.info('checking if remove is needed')
        if os.path.isfile(self.to_path):
            self.log.info('removing')
            os.remove(self.to_path)
