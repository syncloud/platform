from syncloud_app import logger
import getpass

import dns
import port_drill
from port_config import PortConfig
from service_config import ServiceConfig
from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig, PLATFORM_APP_NAME
from syncloud_platform.insider.config import RedirectConfig
from syncloud_platform.insider.port_prober import PortProber
from syncloud_platform.insider.util import protocol_to_port
from syncloud_platform.tools.app import get_app_data_root
from syncloud_platform.tools import network
from syncloud_platform.tools.chown import chown


class Insider:

    def __init__(self, dns_service, platform_config, user_platform_config, port_config, redirect_config):
        self.port_config = port_config
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.dns = dns_service
        self.redirect_config = redirect_config
        self.logger = logger.get_logger('insider')

    def sync_all(self):
        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        self.dns.sync(drill)
        if not getpass.getuser() == self.platform_config.cron_user():
            chown(self.platform_config.cron_user(), self.platform_config.data_dir())

    def acquire_domain(self, email, password, user_domain):
        result = self.dns.acquire(email, password, user_domain)
        return result

    def add_main_device_service(self, protocol, external_access):
        drill = self.get_drill(external_access)
        self.dns.remove_service("server", drill)
        self.dns.add_service("server", protocol, "server", protocol_to_port(protocol), drill)
        self.dns.sync(drill)
        self.user_platform_config.set_protocol(protocol)
        self.user_platform_config.set_external_access(external_access)

    def get_drill(self, external_access):
        drill = port_drill.NonePortDrill()
        if external_access:
            mapper = port_drill.provide_mapper()
            if mapper:
                prober = PortProber(self.user_platform_config, self.user_platform_config.get_redirect_api_url())
                drill = port_drill.PortDrill(self.port_config, mapper, prober)
        return drill


def get_insider():

    data_root = get_app_data_root(PLATFORM_APP_NAME)

    redirect_config = RedirectConfig(data_root)
    platform_config = PlatformConfig()
    user_platform_config = PlatformUserConfig(platform_config.get_user_config())
    port_config = PortConfig(data_root)
    service_config = ServiceConfig(data_root)

    dns_service = dns.Dns(
        service_config,
        network.local_ip(),
        redirect_config,
        user_platform_config)

    return Insider(
        dns_service,
        platform_config,
        user_platform_config,
        port_config,
        redirect_config)