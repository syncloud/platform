import os
from os.path import join
from syncloud.config.config import PlatformConfig, PLATFORM_CONFIG_DIR
from syncloud.tools.app import get_app_data_root
from syncloud.tools.facade import Facade

from port_config import PortConfig
from service_config import ServiceConfig

from syncloud.insider.config import InsiderConfig, RedirectConfig
from syncloud.insider import config
import upnpc
import upnpc_mock
import port_mapper
import dns
import cron


class Insider:

    def __init__(self, port_mapper, dns, cron, insider_config, service_config):
        self.insider_config = insider_config
        self.port_mapper = port_mapper
        self.dns = dns
        self.cron = cron
        self.service_config = service_config

    def list_ports(self):
        return self.port_mapper.list()

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
        return self.port_mapper.get(port)

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

    def set_redirect_info(self, domain, api_url):
        return self.insider_config.update(domain, api_url)

    def endpoints(self):
        return self.dns.endpoints()


def get_insider(config_path=PLATFORM_CONFIG_DIR, use_upnpc_mock=False, data_root=None):

    if not data_root:
        data_root = get_app_data_root('platform')

    redirect_config = RedirectConfig(join(data_root, 'redirect.cfg'))
    insider_config = InsiderConfig(config_path, redirect_config)

    local_ip = Facade().local_ip()

    if use_upnpc_mock or insider_config.is_upnpc_mock():
        upnpclient = upnpc_mock.Upnpc()
    else:
        upnpclient = upnpc.Upnpc(local_ip)

    mapper = port_mapper.PortMapper(
        PortConfig(join(data_root, 'ports.json')),
        upnpclient)

    service_config = ServiceConfig(join(data_root, 'services.json'))

    dns_service = dns.Dns(
        insider_config,
        config.DomainConfig(join(data_root, 'domain.json')),
        service_config,
        mapper,
        local_ip)
    platform_config = PlatformConfig(config_path)
    cron_service = cron.Cron(
        join(platform_config.bin_dir(), 'insider'),
        join(platform_config.data_dir(), 'insider-cron.log'),
        insider_config.get_cron_period_mins())

    return Insider(mapper, dns_service, cron_service, insider_config, service_config)
