import os
from os.path import join
from string import Template
from syncloud_app import logger
from syncloud_platform.gaplib import gen


def proxy_definition(app, port, template_dir, template, www_root_public):
    return Template(open(join(template_dir, template)).read().strip()).substitute(
        {'app': app, 'port': port, 'www_root_public': www_root_public})


class Nginx:
    def __init__(self, platform_config, systemctl, device_info):
        self.systemctl = systemctl
        self.config = platform_config
        self.device_info = device_info
        self.log = logger.get_logger('nginx')

    def add_app(self, app, port):

        self.remove_app(app, reload=False)

        webapp = self.__app_file(app)
        self.log.info('creating {0}'.format(webapp))
        with open(webapp, 'w') as f:
            f.write(proxy_definition(
                app, port, self.config.nginx_config_dir(), 'app.server', self.config.www_root_public()))

        self.systemctl.reload_service('platform.nginx-public')

    def remove_app(self, app, reload=True):

        webapp = self.__app_file(app)
        if os.path.isfile(webapp):
            self.log.info('removing {0}'.format(webapp))
            os.remove(webapp)

        if reload:
            self.systemctl.reload_service('platform.nginx-public')

    def __app_file(self, app):
        return join(self.config.nginx_webapps(), '{0}.server'.format(app))

    def reload_internal(self):
        self.systemctl.reload_service('platform.nginx-internal')
    
    def reload_public(self):
        self.systemctl.reload_service('platform.nginx-public')

    def init_config(self):
        domain = self.device_info.domain()
        nginx_public_template = join(self.config.config_dir(), 'nginx', 'public.conf')
        nginx_public_runtime = join(self.config.data_dir(), 'config.runtime', 'nginx', 'public.conf')
        variables = { 'user_domain': domain }
        gen.generate_file_jinja(nginx_public_template, nginx_public_runtime, variables)