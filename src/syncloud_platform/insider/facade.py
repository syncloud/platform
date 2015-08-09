import os
from os.path import join
from syncloud_platform.config.config import PlatformConfig, PLATFORM_CONFIG_DIR
from syncloud_platform.tools.app import get_app_data_root
from syncloud_platform.tools.facade import Facade

from port_config import PortConfig
from service_config import ServiceConfig

from syncloud_platform.insider.config import InsiderConfig, RedirectConfig, DomainConfig
import port_drill
import dns
import cron


class Insider:

    def __init__(self, mapper, dns, cron, insider_config, service_config):
        self.insider_config = insider_config
        self.mapper = mapper
        self.dns = dns
        self.cron = cron
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
        self.cron.on()
        return result

    def drop_domain(self):
        self.dns.drop()
        self.cron.off()

    def full_name(self):
        return self.dns.full_name()

    def user_domain(self):
        return self.dns.user_domain()

    def cron_on(self):
        return self.cron.on()

    def cron_off(self):
        return self.cron.off()

    def endpoints(self):
        return self.dns.endpoints()


def get_insider(config_path=PLATFORM_CONFIG_DIR):

    data_root = get_app_data_root('platform')

    redirect_config = RedirectConfig(data_root)
    insider_config = InsiderConfig(config_path)

    local_ip = Facade().local_ip()

    port_config = PortConfig(join(data_root, 'ports.json'))

    drill = port_drill.NonePortDrill()
    if insider_config.get_external_access():
        mapper = port_drill.provide_mapper()
        if mapper:
            drill = port_drill.PortDrill(port_config, mapper)

    service_config = ServiceConfig(join(data_root, 'services.json'))

    dns_service = dns.Dns(
        insider_config,
        DomainConfig(join(data_root, 'domain.json')),
        service_config,
        drill,
        local_ip,
        redirect_config)
    platform_config = PlatformConfig(config_path)
    cron_service = cron.Cron(
        join(platform_config.bin_dir(), 'insider'),
        join(platform_config.data_dir(), 'insider-cron.log'),
        insider_config.get_cron_period_mins())

    return Insider(drill, dns_service, cron_service, insider_config, service_config)
