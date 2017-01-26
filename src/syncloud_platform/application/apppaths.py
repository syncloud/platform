from os import mkdir
from os.path import join, isdir
from shutil import rmtree


class AppPaths:

    def __init__(self, app_name, platform_config):
        self.app_name = app_name
        self.platform_config = platform_config

    def get_install_dir(self):
        return join(self.platform_config.apps_root(), self.app_name)

    def get_data_dir(self, remove_existing=False):
        config = self.platform_config
        if not isdir(config.data_root()):
            print("creating app data root: {0}".format(config.data_root()))
            mkdir(config.data_root())

        app_data_dir = join(config.data_root(), self.app_name)
        print("checking app config folder: {0}".format(app_data_dir))

        if isdir(app_data_dir) and remove_existing:
            print("removing existing app data dir: {0}".format(app_data_dir))
            rmtree(app_data_dir, ignore_errors=True)

        if not isdir(app_data_dir):
            print("creating app data dir: {0}".format(app_data_dir))
            mkdir(app_data_dir)
        else:
            print("app data dir exists: {0}".format(app_data_dir))

        return app_data_dir