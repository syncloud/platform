from os.path import dirname, join
import tempfile
import os
from syncloud_platform.config.config import PLATFORM_CONFIG_NAME, PlatformConfig
from syncloud_platform.config.user_config import PlatformUserConfig
from syncloud_platform.insider.port_config import PortConfig, PORT_CONFIG_NAME


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


def get_port_config(mappings):
    config = PortConfig(dirname(temp_file(filename=PORT_CONFIG_NAME)))
    config.save(mappings)
    return config


test_conf_dir = join(dirname(__file__), 'conf')
test_services_config_file_name = 'services.json'
test_services_config_file = join(test_conf_dir, test_services_config_file_name)

CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

platform_config_file = CONFIG_DIR


def get_user_platform_config():
    config = PlatformUserConfig(config_file=temp_file())
    return config


def get_platform_config():
    config = PlatformConfig(platform_config_file)
    return config
