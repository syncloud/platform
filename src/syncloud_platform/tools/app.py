from grp import getgrnam
import os
from os.path import join, isdir
from pwd import getpwnam
import shutil
from syncloud_platform.config.config import PlatformConfig


def get_app_dir(app_name):
    config = PlatformConfig()
    return join(config.apps_root(), app_name)

def get_app_data_dir(app_name, remove_existing=False):
    config = PlatformConfig()
    if not isdir(config.data_root()):
        print("creating app data root: {0}".format(config.data_root()))
        os.mkdir(config.data_root())

    app_data_dir = join(config.data_root(), app_name)
    print("checking app config folder: {0}".format(app_data_dir))

    if isdir(app_data_dir) and remove_existing:
        print("removing existing app data dir: {0}".format(app_data_dir))
        shutil.rmtree(app_data_dir, ignore_errors=True)

    if not isdir(app_data_dir):
        print("creating app data dir: {0}".format(app_data_dir))
        os.mkdir(app_data_dir)
    else:
        print("app data dir exists: {0}".format(app_data_dir))

    return app_data_dir