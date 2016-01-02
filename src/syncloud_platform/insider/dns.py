import getpass
from urlparse import urljoin

import requests
import convertible
from IPy import IP
from syncloud_app import logger
from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.insider import util
from syncloud_platform.insider.config import Service
from syncloud_platform.tools import id
from syncloud_platform.tools.chown import chown


class ServiceUrls:
    def __init__(self, device_domain, service, mapping):
        self.domain_port = "{}:{}".format(device_domain, mapping.external_port)
        self.full_url = "{}://{}/{}".format(service.protocol, self.domain_port, service.url)


class Endpoint:
    def __init__(self, service, domain, external_port):
        self.service = service
        self.external_host = domain
        self.external_port = external_port


class Dns:

    def __init__(self, domain_config, service_config, local_ip, redirect_config, platform_config=None, fix_permissions=True):
        self.fix_permissions = fix_permissions
        self.redirect_config = redirect_config
        self.local_ip = local_ip
        self.domain_config = domain_config
        self.service_config = service_config
        self.logger = logger.get_logger('dns')
        if platform_config:
            self.config = platform_config
        else:
            self.config = PlatformConfig()

    def acquire(self, email, password, user_domain):
        device_id = id.id()
        data = {
            'email': email,
            'password': password,
            'user_domain': user_domain,
            'device_mac_address': device_id.mac_address,
            'device_name': device_id.name,
            'device_title': device_id.title,
        }
        url = urljoin(self.redirect_config.get_api_url(), "/domain/acquire")
        response = requests.post(url, data)
        util.check_http_error(response)
        domain = convertible.from_json(response.text)
        self.domain_config.save(domain)
        return domain

    def add_service(self, name, protocol, service_type, port, port_drill):
        port_drill.sync_new_port(port)
        new_service = Service(name, protocol, service_type, port)
        self.service_config.add_or_update(new_service)

    def get_service(self, name):
        return self.service_config.get(name)

    def get_service_by_port(self, port):
        return self.service_config.get(port)

    def remove_service(self, name, port_drill):
        service = self.get_service(name)
        if service:
            self.service_config.remove(name)
            port_drill.remove(service.port)

    def full_name(self):
        return '{}.{}'.format(self.user_domain(), self.redirect_config.get_domain())

    def user_domain(self):
        return self.domain_config.load().user_domain

    def sync(self, port_drill):
        port_drill.sync()

        services = self.service_config.load()
        if not self.domain_config.exists():
            self.logger.info('nothing to sync yet, no dns configuration')
            return

        domain = self.domain_config.load()
        if not domain:
            raise Exception("No token saved, need to call set_dns or get_dns_token first")

        services_data = []
        for service in services:
            self.logger.debug('service: {0} '.format(service.port))
            mapping = port_drill.get(service.port)
            if mapping:
                service.port = mapping.external_port
                service_data = dict(
                    name=service.name,
                    protocol=service.protocol,
                    type=service.type,
                    url=service.url,
                    port=mapping.external_port,
                    local_port=mapping.local_port
                )
                services_data.append(service_data)

        data = {
            'token': domain.update_token,
            'local_ip': self.local_ip,
            'services': services_data,
            'map_local_address': False}

        external_ip = port_drill.external_ip()

        if not external_ip:
            self.logger.warn("No external ip")
        else:
            if IP(external_ip).iptype() != 'PUBLIC':
                external_ip = None
                self.logger.warn("External ip is not public")

        if external_ip:
            data['ip'] = external_ip
        else:
            data['map_local_address'] = True
            self.logger.warn("Will try server side client ip detection")

        url = urljoin(self.redirect_config.get_api_url(), "/domain/update")

        self.logger.debug('url: ' + url)
        json = convertible.to_json(data)
        self.logger.debug('request: ' + json)
        response = requests.post(url, json)

        util.check_http_error(response)

        if not getpass.getuser() == self.config.cron_user() and self.fix_permissions:
            chown(self.config.cron_user(), self.config.data_dir())

