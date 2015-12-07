from syncloud_app import logger

import cron
import dns
import port_drill
from port_config import PortConfig
from service_config import ServiceConfig
from syncloud_platform.config.config import PlatformConfig, PLATFORM_CONFIG_DIR, PlatformUserConfig, PLATFORM_APP_NAME
from syncloud_platform.insider.config import RedirectConfig, DomainConfig
from syncloud_platform.insider.port_prober import PortProber
from syncloud_platform.insider.util import protocol_to_port
from syncloud_platform.tools.app import get_app_data_root
from syncloud_platform.tools.events import trigger_app_event_domain
from syncloud_platform.tools import network


class Insider:

    def __init__(self, dns_service, platform_cron, user_platform_config, port_config, platform_config):
        self.platform_config = platform_config
        self.port_config = port_config
        self.user_platform_config = user_platform_config
        self.dns = dns_service
        self.platform_cron = platform_cron
        self.logger = logger.get_logger('insider')

    def sync_all(self):
        return self.dns.sync()

    def acquire_domain(self, email, password, user_domain):
        result = self.dns.acquire(email, password, user_domain)
        self.sync_all()
        self.platform_cron.remove()
        self.platform_cron.create()
        return result

    def add_main_device_service(self, mode='http'):
        drill = get_drill(True, self.dns.domain_config, self.port_config, self.dns.redirect_config)
        self.dns.remove_service("server", drill)
        self.dns.add_service("server", mode, "server", protocol_to_port(mode), drill)
        self.sync_all()
        self.user_platform_config.set_external_access(mode)
        trigger_app_event_domain(self.platform_config.apps_root())

    def remove_main_device_service(self):
        drill = get_drill(True, self.dns.domain_config, self.port_config, self.dns.redirect_config)
        self.dns.remove_service("server", drill)
        self.sync_all()
        self.user_platform_config.disable_external_access()
        trigger_app_event_domain(self.platform_config.apps_root())


def get_insider(config_path=PLATFORM_CONFIG_DIR):

    data_root = get_app_data_root(PLATFORM_APP_NAME)

    redirect_config = RedirectConfig(data_root)
    user_platform_config = PlatformUserConfig()
    port_config = PortConfig(data_root)
    domain_config = DomainConfig(data_root)

    drill = get_drill(user_platform_config.get_external_access(), domain_config, port_config, redirect_config)

    dns_service = dns.Dns(
        domain_config,
        ServiceConfig(data_root),
        drill,
        network.local_ip(),
        redirect_config)

    platform_config = PlatformConfig(config_path)

    return Insider(
        dns_service,
        cron.PlatformCron(platform_config),
        user_platform_config,
        port_config,
        platform_config)


def get_drill(external_access, domain_config, port_config, redirect_config):

    drill = port_drill.NonePortDrill()
    if external_access:
        mapper = port_drill.provide_mapper()
        if mapper:
            prober = PortProber(domain_config, redirect_config.get_api_url())
            drill = port_drill.PortDrill(port_config, mapper, prober)
    return drill
