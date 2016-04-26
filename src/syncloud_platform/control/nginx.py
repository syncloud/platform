import os
from os.path import join
from string import Template
from syncloud_platform.systemd.systemctl import reload_service
from syncloud_app import logger


class Nginx:
    def __init__(self, platform_config):
        self.config = platform_config
        self.log = logger.get_logger('nginx')

    def add_app(self, app, port):

        self.remove_app(app, reload=False)

        webapp = self.__app_file(app)
        self.log.info('creating {0}'.format(webapp))
        with open(webapp, 'w') as f:
            f.write(self.proxy_definition(app, port, self.config.nginx_config_dir(), 'app.server'))

        reload_service('platform-nginx')

    def proxy_definition(self, app, port, template_dir, template):
        return Template(open(join(template_dir, template)).read().strip()).substitute({'app': app, 'port': port})

    def remove_app(self, app, reload=True):

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            self.log.info('removing {0}'.format(webapp))
            os.remove(webapp)

        if reload:
            reload_service('platform-nginx')

    def __app_file(self, app):
        return join(self.config.nginx_webapps(), '{0}.server'.format(app))

    def reload(self):
        reload_service('platform-nginx')

