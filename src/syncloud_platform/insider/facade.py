import getpass

from syncloud_app import logger
from syncloud_platform.insider.dns import Dns
from syncloud_platform.insider.port_drill import PortDrillFactory
from syncloud_platform.insider.port_config import PortConfig
from syncloud_platform.insider.service_config import ServiceConfig
from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig, PLATFORM_APP_NAME
from syncloud_platform.insider.util import protocol_to_port
from syncloud_platform.tools.app import get_app_data_root
from syncloud_platform.tools import network
from syncloud_platform.tools.chown import chown


class Insider:

    def __init__(self, dns_service, platform_config, user_platform_config, port_drill_factory):
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.dns_service = dns_service
        self.port_drill_factory = port_drill_factory
        self.logger = logger.get_logger('insider')

    def sync_all(self):
        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        self.dns_service.sync(drill)
        if not getpass.getuser() == self.platform_config.cron_user():
            chown(self.platform_config.cron_user(), self.platform_config.data_dir())

    def acquire_domain(self, email, password, user_domain):
        result = self.dns_service.acquire(email, password, user_domain)
        return result

    def add_main_device_service(self, protocol, external_access):
        drill = self.get_drill(external_access)
        self.dns_service.remove_service("server", drill)
        self.dns_service.add_service("server", protocol, "server", protocol_to_port(protocol), drill)
        self.dns_service.sync(drill)
        self.user_platform_config.update_device_access(external_access, protocol)

    def get_drill(self, external_access):
        return self.port_drill_factory.get_drill(external_access)


def get_insider():

    data_root = get_app_data_root(PLATFORM_APP_NAME)

    platform_config = PlatformConfig()
    user_platform_config = PlatformUserConfig(platform_config.get_user_config())
    port_config = PortConfig(data_root)
    service_config = ServiceConfig(data_root)

    dns_service = Dns(
        service_config,
        network.local_ip(),
        user_platform_config)

    port_drill_factory = PortDrillFactory(user_platform_config, port_config)

    return Insider(
        dns_service,
        platform_config,
        user_platform_config,
        port_drill_factory)