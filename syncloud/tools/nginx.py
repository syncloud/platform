import os
from os.path import join
from syncloud.systemd.systemctl import reload_service
from syncloud.config.config import PlatformConfig


class Nginx:
    def __init__(self):
        self.config = PlatformConfig()

    def add_app(self, app, port):

        self.remove_app(app, reload=False)

        proxy = '''location /${0}/ {
            proxy_pass http://127.0.0.1:${1}/${0}/;
        }'''.format(app, port)

        webapp = self.__app_file(app)
        with open(webapp, 'w') as f:
            f.write(proxy)

        reload_service('platform-nginx')

    def remove_app(self, app, reload=True):

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            os.remove(webapp)
        if reload:
            reload_service('platform-nginx')

    def __app_file(self, app):
        return join(self.config.nginx_webapps(), app)
