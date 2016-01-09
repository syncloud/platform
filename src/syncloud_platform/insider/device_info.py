from syncloud_platform.insider.util import protocol_to_port


def construct_url(protocol, external_port, domain, app=None):
    protocol = protocol or 'http'
    if external_port in [80, 443]:
        external_port = ''
    else:
        external_port = ':{0}'.format(external_port)
    app_string = ''
    if app:
        app_string = app + '.'
    return '{0}://{1}{2}{3}'.format(protocol, app_string, domain, external_port)


class DeviceInfo:
    def __init__(self, user_platform_config, port_config):
        self.port_config = port_config
        self.user_platform_config = user_platform_config

    def domain(self):
        user_domain = self.user_platform_config.get_user_domain()
        if user_domain is not None:
            return '{0}.{1}'.format(user_domain, self.user_platform_config.get_redirect_domain())
        else:
            return None

    def url(self, app=None):
        protocol = self.user_platform_config.get_protocol()
        port = protocol_to_port(protocol)
        if self.user_platform_config.get_external_access():
            mapping = self.port_config.get(port)
            if mapping:
                port = mapping.external_port

        domain_name = self.domain()
        if domain_name:
            return construct_url(protocol, port, domain_name, app)
        else:
            ''
