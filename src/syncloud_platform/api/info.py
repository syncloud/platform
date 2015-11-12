from syncloud_platform.insider.config import DomainConfig, RedirectConfig


def domain():
    domain_config = DomainConfig()
    redirect = RedirectConfig()
    return '{0}.{1}'.format(domain_config.load().user_domain, redirect.get_domain())
