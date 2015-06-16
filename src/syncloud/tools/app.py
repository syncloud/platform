from grp import getgrnam
import os
from os.path import join
from pwd import getpwnam
import shutil
from syncloud.config.config import PlatformConfig

def get_app_data_root(app_name, user=None):
    config = PlatformConfig()
    if not os.path.isdir(config.data_root()):
        print("creating app data root: {0}".format(config.data_root()))
        os.mkdir(config.data_root())

    return create_data_dir(config.data_root(), app_name, user)

def create_data_dir(app_data_dir, dir_name, user=None, remove_existing=False):
    data_dir = join(app_data_dir, dir_name)
    print("checking app config folder: {0}".format(data_dir))

    if os.path.isdir(data_dir) and remove_existing:
        print("removing existing app data dir: {0}".format(data_dir))
        shutil.rmtree(data_dir, ignore_errors=True)

    if not os.path.isdir(data_dir):
        print("creating app data dir: {0}".format(data_dir))
        os.mkdir(data_dir)
        if user:
            print("setting permissions for {0} to {1}".format(data_dir, user))
            os.chown(data_dir, getpwnam(user).pw_uid, getgrnam(user).gr_gid)
    else:
        print("app data dir exists: {0}".format(data_dir))

    return data_dir
