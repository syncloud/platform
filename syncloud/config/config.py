from ConfigParser import ConfigParser


class PlatformConfig:

    def __init__(self, filename='/opt/app/platform/config/platform.cfg'):
        self.parser = ConfigParser()
        self.parser.read(filename)

    def apps_root(self):
        return self.parser.get('platform', 'apps_root')

    def data_root(self):
        return self.parser.get('platform', 'data_root')

    def nginx_webapps(self):
        return self.parser.get('platform', 'nginx_webapps')
