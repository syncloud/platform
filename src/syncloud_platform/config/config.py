from ConfigParser import ConfigParser
from os.path import isfile, join

PLATFORM_CONFIG_DIR = '/opt/app/platform/config'
PLATFORM_CONFIG_NAME = 'platform.cfg'


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

    def __get(self, key):
        return self.parser.get('platform', key)

    def __set(self, key, value):
        self.parser.set('platform', key, value)
        with open(self.filename, 'wb') as f:
            self.parser.write(f)


class PlatformUserConfig:

    def __init__(self):
        self.parser = ConfigParser()
        self.filename = PlatformConfig().get_user_config()
        if not isfile(self.filename):
            self.parser.add_section('platform')
            self.set_activated(False)
            self.__save()
        else:
            self.parser.read(self.filename)

    def is_activated(self):
        return self.parser.get('platform', 'activated')

    def set_activated(self, value):
        self.parser.set('platform', 'activated', value)
        self.__save()

    def __save(self):
        with open(self.filename, 'wb') as f:
            self.parser.write(f)
