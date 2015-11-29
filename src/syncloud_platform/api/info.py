from syncloud_platform.config.config import PlatformUserConfig
from syncloud_platform.insider.config import DomainConfig, RedirectConfig
from syncloud_platform.insider.port_config import PortConfig


def domain():
    domain_config = DomainConfig()
    redirect = RedirectConfig()
    return '{0}.{1}'.format(domain_config.load().user_domain, redirect.get_domain())


def url():
    external_access_protocol = PlatformUserConfig().get_external_access()
    port = 80
    if external_access_protocol:
        port = PortConfig().get(port).external_port

    return __url(external_access_protocol, port, domain())


def __url(protocol, external_port, domain):
    protocol = protocol or 'http'
    if external_port in [80, 443]:
        external_port = ''
    else:
        external_port = ':{0}'.format(external_port)
    return '{0}://{1}{2}'.format(protocol, domain, external_port)
