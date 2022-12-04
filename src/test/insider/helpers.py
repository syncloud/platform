import tempfile
import os
from syncloud_platform.config.user_config import PlatformUserConfig


def temp_file(text='', filename=None):
    if filename:
        filename = '/tmp/' + filename
        with open(filename, 'w') as f:
            f.write(text)
    else:
        fd, filename = tempfile.mkstemp()
        f = os.fdopen(fd, 'w')
        f.write(text)
        f.close()
    return filename


def get_user_platform_config():
    config = PlatformUserConfig(temp_file())
    config.init_user_config()
    return config
