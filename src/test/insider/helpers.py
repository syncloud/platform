from os.path import dirname, join
import tempfile
import os
from syncloud_platform.config.config import PlatformConfig, PLATFORM_CONFIG_NAME

from syncloud_platform.insider.port_config import PortConfig
from syncloud_platform.insider.service_config import ServiceConfig

from syncloud_platform.insider.config import DomainConfig, InsiderConfig, RedirectConfig, INSIDER_CONFIG_NAME


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
    filename = temp_file()
    config = PortConfig(filename)
    config.save(mappings)
    return config


def get_domain_config(domain=None):
    domain_file = temp_file()
    domain_config = DomainConfig(domain_file)
    domain_config.save(domain)
    return domain_config


def get_service_config(services):
    filename = temp_file()
    config = ServiceConfig(filename)
    config.save(services)
    return config


test_conf_dir = join(dirname(__file__), 'conf')
test_services_config_file_name = 'services.json'
test_services_config_file = join(test_conf_dir, test_services_config_file_name)

CONFIG_DIR = join(dirname(__file__), '..', '..', '..', 'config')

insider_config_file = join(CONFIG_DIR, INSIDER_CONFIG_NAME)
insider_config = open(insider_config_file).read()

platform_config_file = join(CONFIG_DIR, PLATFORM_CONFIG_NAME)


def get_redirect_config():
    redirect_config = RedirectConfig(temp_file())
    redirect_config.update('domain.com', 'http://api.domain.com')
    return redirect_config


def get_insider_config():
    tem_config = temp_file(insider_config, INSIDER_CONFIG_NAME)
    temp_config_dir = dirname(tem_config)
    config = InsiderConfig(temp_config_dir)
    return config

