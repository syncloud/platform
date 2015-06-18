import os
from os.path import join
from string import Template
from syncloud.systemd.systemctl import reload_service
from syncloud.config.config import PlatformConfig


class Nginx:
    def __init__(self):
        self.config = PlatformConfig()

    def add_app(self, app, port):

        self.remove_app(app, reload=False)

        with open(self.__app_file(app), 'w') as f:
            f.write(self.proxy_definition(app, port))

        reload_service('platform-nginx')

    def proxy_definition(self, app, port):
        return Template(
            'location /${app} {\n'
            '    proxy_pass      http://localhost:${port}/${app}/ ;\n'
            '    proxy_redirect  http://localhost:${port}/${app}/ /${app}/ ;\n'
            '}').substitute(
            {'app': app, 'port': port}
        )

    def remove_app(self, app, reload=True):

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            os.remove(webapp)
        if reload:
            reload_service('platform-nginx')

    def __app_file(self, app):
        return join(self.config.nginx_webapps(), app)
