from os.path import join
from syncloud_platform.config.config import PlatformConfig, PLATFORM_CONFIG_DIR, PlatformUserConfig
from syncloud_platform.insider.port_prober import PortProber
from syncloud_platform.insider.util import port_to_protocol, protocol_to_port
from syncloud_platform.tools.app import get_app_data_root
from syncloud_platform.tools.facade import Facade

from port_config import PortConfig
from service_config import ServiceConfig

from syncloud_platform.insider.config import RedirectConfig, DomainConfig
import port_drill
import dns
import cron


class Insider:

    def __init__(self, mapper, dns, platform_cron, service_config, user_platform_config):
        self.user_platform_config = user_platform_config
        self.mapper = mapper
        self.dns = dns
        self.platform_cron = platform_cron
        self.service_config = service_config

    def list_ports(self):
        return self.mapper.list()

    def sync_all(self):
        return self.dns.sync()

    def add_service(self, name, protocol, type, port, url):
        result = self.dns.add_service(name, protocol, type, port, url)
        self.sync_all()
        return result

    def remove_service(self, name):
        result = self.dns.remove_service(name)
        self.sync_all()
        return result

    def get_service(self, name):
        return self.dns.get_service(name)

    def get_mapping(self, port):
        return self.mapper.get(port)

    def service_info(self, name):
        return self.dns.service_info(name)

    def acquire_domain(self, email, password, user_domain):
        result = self.dns.acquire(email, password, user_domain)
        self.sync_all()
        self.platform_cron.remove()
        self.platform_cron.create()
        return result

    def drop_domain(self):
        self.dns.drop()
        self.platform_cron.remove()

    def full_name(self):
        return self.dns.full_name()

    def user_domain(self):
        return self.dns.user_domain()

    def cron_on(self):
        return self.platform_cron.create()

    def cron_off(self):
        return self.platform_cron.remove()

    def endpoints(self):
        return self.dns.endpoints()

    def add_main_device_service(self, mode='http'):
        self.add_service("server", mode, "server", protocol_to_port(mode), None)
        self.user_platform_config.set_external_access(mode)

    def remove_main_device_service(self):
        self.remove_service("server")

def get_insider(config_path=PLATFORM_CONFIG_DIR):

    data_root = get_app_data_root('platform')

    redirect_config = RedirectConfig(data_root)
    user_platform_config = PlatformUserConfig()
    local_ip = Facade().local_ip()

    port_config = PortConfig(join(data_root, 'ports.json'))
    domain_config = DomainConfig(join(data_root, 'domain.json'))

    drill = port_drill.NonePortDrill()
    if user_platform_config.get_external_access():
        mapper = port_drill.provide_mapper()
        if mapper:
            prober = PortProber(domain_config, redirect_config.get_api_url())
            drill = port_drill.PortDrill(port_config, mapper, prober)

    service_config = ServiceConfig(join(data_root, 'services.json'))

    dns_service = dns.Dns(
        domain_config,
        service_config,
        drill,
        local_ip,
        redirect_config)

    return Insider(
        drill,
        dns_service,
        cron.PlatformCron(PlatformConfig(config_path)),
        service_config,
        user_platform_config)
