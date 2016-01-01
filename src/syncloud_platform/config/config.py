from ConfigParser import ConfigParser
from os.path import isfile, join
from syncloud_app import logger

PLATFORM_CONFIG_DIR = '/opt/app/platform/config'
PLATFORM_CONFIG_NAME = 'platform.cfg'
PLATFORM_APP_NAME = 'platform'


class PlatformConfig:

    def __init__(self, config_dir=PLATFORM_CONFIG_DIR):
        self.parser = ConfigParser()
        self.filename = join(config_dir, PLATFORM_CONFIG_NAME)
        self.parser.read(self.filename)

    def apps_root(self):
        return self.__get('apps_root')

    def data_root(self):
        return self.__get('data_root')

    def www_root(self):
        return self.__get('www_root')

    def app_dir(self):
        return self.__get('app_dir')

    def data_dir(self):
        return self.__get('data_dir')

    def config_dir(self):
        return self.__get('config_dir')

    def bin_dir(self):
        return self.__get('bin_dir')

    def nginx_webapps(self):
        return self.__get('nginx_webapps')

    def nginx_config_dir(self):
        return self.__get('nginx_config_dir')

    def cron_user(self):
        return self.__get('cron_user')

    def cron_cmd(self):
        return self.__get('cron_cmd')

    def cron_schedule(self):
        return self.__get('cron_schedule')

    def get_web_secret_key(self):
        return self.__get('web_secret_key')

    def set_web_secret_key(self, value):
        return self.__set('web_secret_key', value)

    def get_user_config(self):
        return self.__get('user_config')

    def get_log_root(self):
        return self.__get('log_root')

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

    def set_ssh_port(self, value):
        return self.__set('ssh_port', value)

    def get_rest_internal_log(self):
        return self.__get('rest_internal_log')

    def get_rest_public_log(self):
        return self.__get('rest_public_log')

    def get_ssl_certificate_file(self):
        return self.__get('ssl_certificate_file')

    def get_default_ssl_certificate_file(self):
        return self.__get('default_ssl_certificate_file')

    def get_ssl_key_file(self):
        return self.__get('ssl_key_file')

    def get_default_ssl_key_file(self):
        return self.__get('default_ssl_key_file')

    def get_openssl_config(self):
        return self.__get('openssl_config')

    def get_platform_log(self):
        return self.__get('platform_log')

    def __get(self, key):
        return self.parser.get('platform', key)

    def __set(self, key, value):
        self.parser.set('platform', key, value)
        with open(self.filename, 'wb') as f:
            self.parser.write(f)


class PlatformUserConfig:

    def __init__(self, config_file=None):
        self.logger = logger.get_logger('PlatformUserConfig')
        self.parser = ConfigParser()
        if config_file:
            self.filename = config_file
        else:
            self.filename = PlatformConfig().get_user_config()

        if not isfile(self.filename):
            self.parser.add_section('platform')
            self.set_activated(False)
            self.__save()
        else:
            self.parser.read(self.filename)

        if not self.parser.has_section('platform'):
            self.parser.add_section('platform')

    def is_activated(self):
        return self.parser.getboolean('platform', 'activated')

    def set_activated(self, value):
        self.parser.set('platform', 'activated', str(value))
        self.__save()

    def get_external_access(self):
        self.parser.read(self.filename)
        external_access = False
        if self.parser.has_option('platform', 'external_access'):
            external_access = self.parser.getboolean('platform', 'external_access')
        self.logger.info('external_access = {0}'.format(external_access))
        return external_access

    def set_external_access(self, external_access):
        self.logger.info('saving config, set_external_access = {0}'.format(external_access))
        self.parser.read(self.filename)
        self.parser.set('platform', 'external_access', external_access)
        self.__save()

    def get_protocol(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'protocol'):
            return 'http'
        return self.parser.get('platform', 'protocol')

    def set_protocol(self, protocol):
        self.parser.read(self.filename)
        self.parser.set('platform', 'protocol', protocol)
        self.__save()

    def __save(self):
        with open(self.filename, 'wb') as f:
            self.parser.write(f)

