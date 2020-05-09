import os
from configparser import ConfigParser
from os.path import isfile
from os.path import join

PLATFORM_CONFIG_NAME = 'platform.cfg'
PLATFORM_APP_NAME = 'platform'
WEB_CERTIFICATE_PORT = 80
WEB_ACCESS_PORT = 443
WEB_PROTOCOL = 'https'

APPS_ROOT = '/snap'
DATA_ROOT = '/var/snap'

def env(key, default_value):
    if key in os.environ:
        return os.environ[key]
    return default_value

INSTALL_DIR = env('SNAP', 'not_set')
DATA_DIR = env('SNAP_COMMON', 'not_set')
APP_DATA_PREFIX = 'common/'


class PlatformConfig:

    def __init__(self, config_dir):
        self.parser = ConfigParser()
        self.filename = join(config_dir, PLATFORM_CONFIG_NAME)
        if (not isfile(self.filename)):
            raise Exception('platform config does not exist: {0}'.format(self.filename))
        self.parser.read(self.filename)

    def apps_root(self):
        return self.__get('apps_root')

    def data_root(self):
        return self.__get('data_root')

    def configs_root(self):
        return self.__get('configs_root')

    def config_root(self):
        return self.__get('config_root')

    def www_root_public(self):
        return self.__get('www_root_public')

    def app_dir(self):
        return self.__get('app_dir')

    def data_dir(self):
        return self.__get('data_dir')

    def config_dir(self):
        return self.__get('config_dir')

    def bin_dir(self):
        return self.__get('bin_dir')

    def nginx_config_dir(self):
        return self.__get('nginx_config_dir')

    def cron_user(self):
        return self.__get('cron_user')

    def cron_cmd(self):
        return self.__get('cron_cmd')

    def openssl(self):
        return self.__get('openssl')

    def nginx(self):
        return self.__get('nginx')

    def cron_schedule(self):
        return self.__get('cron_schedule')

    def get_log_root(self):
        return self.__get('log_root')

    def get_log_sender_pattern(self):
        return self.__get('log_sender_pattern')

    def get_internal_disk_dir(self):
        return self.__get('internal_disk_dir')

    def get_external_disk_dir(self):
        return self.__get('external_disk_dir')

    def get_disk_link(self):
        return self.__get('disk_link')

    def get_disk_root(self):
        return self.__get('disk_root')

    def get_ssh_port(self):
        return self.__get('ssh_port')

    def get_rest_internal_log(self):
        return self.__get('rest_internal_log')

    def get_rest_public_log(self):
        return self.__get('rest_public_log')

    def get_ssl_certificate_file(self):
        return self.__get('ssl_certificate_file')

    def get_ssl_ca_certificate_file(self):
        return self.__get('ssl_ca_certificate_file')

    def get_ssl_ca_serial_file(self):
        return self.__get('ssl_ca_serial_file')

    def get_ssl_certificate_request_file(self):
        return self.__get('ssl_certificate_request_file')

    def get_default_ssl_certificate_file(self):
        return self.__get('default_ssl_certificate_file')

    def get_ssl_key_file(self):
        return self.__get('ssl_key_file')

    def get_ssl_ca_key_file(self):
        return self.__get('ssl_ca_key_file')

    def get_default_ssl_key_file(self):
        return self.__get('default_ssl_key_file')

    def get_openssl_config(self):
        return self.__get('openssl_config')

    def get_platform_log(self):
        return self.__get('platform_log')

    def get_hooks_root(self):
        return self.__get('hooks_root')

    def is_certbot_test_cert(self):
        return self.parser.getboolean('platform', 'certbot_test_cert')

    def get_channel(self):
        return self.__get('channel')

    def __get(self, key):
        return self.parser.get('platform', key)
