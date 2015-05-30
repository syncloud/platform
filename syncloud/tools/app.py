import os
from os.path import join

APP_DATA_ROOT = '/opt/data'


def get_data_dir(app_name):

    if not os.path.isdir(APP_DATA_ROOT):
        os.mkdir(APP_DATA_ROOT)

    app_data_dir = join(APP_DATA_ROOT, app_name)
    if not os.path.isdir(app_data_dir):
        os.mkdir(app_data_dir)

    return app_data_dir
