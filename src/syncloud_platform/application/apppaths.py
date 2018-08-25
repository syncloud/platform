from os import mkdir
from os.path import join, isdir
from shutil import rmtree


class AppPaths:

    def __init__(self, app_name, platform_config):
        self.app_name = app_name
        self.platform_config = platform_config

    def get_install_dir(self):
        return join(self.platform_config.apps_root(), self.app_name, 'current')

    def get_data_dir(self, remove_existing=False):
        return join(self.platform_config.data_root(), self.app_name, 'common')
