from os import mkdir
from os.path import join, isdir
from shutil import rmtree


class AppPaths:

    def __init__(self, app_name, platform_config):
        self.app_name = app_name
        self.platform_config = platform_config

