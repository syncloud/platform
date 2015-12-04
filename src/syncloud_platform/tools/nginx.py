import os
from os.path import join
from string import Template
from syncloud_platform.systemd.systemctl import reload_service
from syncloud_platform.config.config import PlatformConfig


class Nginx:
    def __init__(self):
        self.config = PlatformConfig()

    def add_app(self, app, port):

        self.remove_app(app, reload=False)

        with open(self.__app_file('{0}.location'.format(app)), 'w') as f:
            f.write(self.proxy_definition(app, port, self.config.nginx_config_dir(), 'app.location'))

        with open(self.__app_file('{0}.server'.format(app)), 'w') as f:
            f.write(self.proxy_definition(app, port, self.config.nginx_config_dir(), 'app.server'))

        reload_service('platform-nginx')

    def proxy_definition(self, app, port, template_dir, template):
        return Template(open(join(template_dir, template)).read().strip()).substitute({'app': app, 'port': port})

    def remove_app(self, app, reload=True):

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            os.remove(webapp)

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            os.remove(webapp)

        if reload:
            reload_service('platform-nginx')

    def __app_file(self, app):
        return join(self.config.nginx_webapps(), app)
