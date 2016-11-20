import os
from os.path import join
from string import Template
from syncloud_app import logger


class Nginx:
    def __init__(self, platform_config, systemctl):
        self.systemctl = systemctl
        self.config = platform_config
        self.log = logger.get_logger('nginx')

    def add_app(self, app, port):

        self.remove_app(app, reload=False)

        webapp = self.__app_file(app)
        self.log.info('creating {0}'.format(webapp))
        with open(webapp, 'w') as f:
            f.write(self.proxy_definition(app, port, self.config.nginx_config_dir(), 'app.server', self.config.www_root_public()))

        self.systemctl.reload_service('platform-nginx')

    def proxy_definition(self, app, port, template_dir, template, www_root_public):
        return Template(open(join(template_dir, template)).read().strip()).substitute({'app': app, 'port': port, 'www_root_public': www_root_public})

    def remove_app(self, app, reload=True):

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            self.log.info('removing {0}'.format(webapp))
            os.remove(webapp)

        if reload:
            self.systemctl.reload_service('platform-nginx')

    def __app_file(self, app):
        return join(self.config.nginx_webapps(), '{0}.server'.format(app))

    def reload(self):
        self.systemctl.reload_service('platform-nginx')

