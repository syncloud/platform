from ConfigParser import ConfigParser


class PlatformConfig:

    def __init__(self, filename='/opt/app/platform/config/platform.cfg'):
        self.parser = ConfigParser()
        self.parser.read(filename)

    def apps_root(self):
        return self.__get('apps_root')

    def data_root(self):
        return self.__get('data_root')

    def www_root(self):
        return self.__get('www_root')

    def app_dir(self):
        return self.__get('app_dir')

    def config_dir(self):
        return self.__get('config_dir')

    def nginx_webapps(self):
        return self.__get('nginx_webapps')

    def __get(self, key):
        return self.parser.get('platform', key)
