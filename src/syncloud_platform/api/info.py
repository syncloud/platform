from syncloud_platform.config.config import PlatformUserConfig
from syncloud_platform.insider.config import DomainConfig, RedirectConfig
from syncloud_platform.insider.port_config import PortConfig
from syncloud_platform.insider.util import protocol_to_port


def domain():
    domain_config = DomainConfig()
    if domain_config.exists():
        redirect = RedirectConfig()
        return '{0}.{1}'.format(domain_config.load().user_domain, redirect.get_domain())
    else:
        return None


def url(app=None):
    config = PlatformUserConfig()
    protocol = config.get_protocol()
    port = protocol_to_port(protocol)
    if config.get_external_access():
        mapping = PortConfig().get(port)
        if mapping:
            port = mapping.external_port

    domain_name = domain()
    if domain_name:
        return __url(protocol, port, domain_name, app)
    else:
        ''


def __url(protocol, external_port, domain, app=None):
    protocol = protocol or 'http'
    if external_port in [80, 443]:
        external_port = ''
    else:
        external_port = ':{0}'.format(external_port)
    app_string = ''
    if app:
        app_string = app + '.'
    return '{0}://{1}{2}{3}'.format(protocol, app_string, domain, external_port)
