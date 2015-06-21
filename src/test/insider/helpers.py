from os.path import dirname, join
import tempfile
import os
from syncloud.config.config import PlatformConfig, PLATFORM_CONFIG_NAME

from syncloud.insider.port_config import PortConfig
from syncloud.insider.service_config import ServiceConfig

from syncloud.insider.config import DomainConfig, InsiderConfig, RedirectConfig, INSIDER_CONFIG_NAME


def temp_file(text=''):
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

def get_insider_config(domain, api_url):
    config = InsiderConfig(temp_file(insider_config), RedirectConfig(temp_file()))
    config.update(domain, api_url)
    return config

