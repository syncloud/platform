from ConfigParser import ConfigParser


class PlatformConfig:

    def __init__(self, filename='/opt/app/platform/config/platform.cfg'):
        self.parser = ConfigParser()
        self.filename = filename
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

    def get_web_secret_key(self):
        return self.__get('web_secret_key')

    def set_web_secret_key(self, value):
        return self.__set('web_secret_key', value)

    def __get(self, key):
        return self.parser.get('platform', key)

    def __set(self, key, value):
        self.parser.set('platform', key, value)
        with open(self.filename, 'wb') as file:
            self.parser.write(file)
