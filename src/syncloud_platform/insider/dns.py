from urlparse import urljoin

import requests
import convertible
from IPy import IP
from syncloud_app import logger
from syncloud_platform.insider.config import Service
from syncloud_platform.tools import id
from syncloud_app.main import PassthroughJsonError


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

    def __init__(self, insider_config, domain_config, service_config, port_mapper, local_ip):
        self.local_ip = local_ip
        self.domain_config = domain_config
        self.service_config = service_config
        self.port_mapper = port_mapper
        self.insider_config = insider_config
        self.logger = logger.get_logger('dns')

    def check_http_error(self, response):
        if not response.status_code == 200:
            self.logger.error(response.text)
            error = convertible.from_json(response.text)
            raise PassthroughJsonError(error.message, response.text)

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
        url = urljoin(self.insider_config.get_redirect_api_url(), "/domain/acquire")
        response = requests.post(url, data)
        self.check_http_error(response)
        domain = convertible.from_json(response.text)
        self.domain_config.save(domain)
        return domain

    def drop(self):
        self.domain_config.remove()
        self.service_config.remove_all()
        self.port_mapper.remove_all()

    def add_service(self, name, protocol, service_type, port, url=None):
        self.port_mapper.sync_new_port(port)
        new_service = Service(name, protocol, service_type, port, url)
        self.service_config.add_or_update(new_service)

    def get_service(self, name):
        return self.service_config.get(name)

    def get_service_by_port(self, port):
        return self.service_config.get(port)

    def remove_service(self, name):
        service = self.get_service(name)
        if service:
            self.service_config.remove(name)
            self.port_mapper.remove(service.port)

    def full_name(self):
        return '{}.{}'.format(self.user_domain(), self.insider_config.get_redirect_main_domain())

    def user_domain(self):
        return self.domain_config.load().user_domain

    def service_to_endpoint(self, service):
        mapping = self.port_mapper.get(service.port)
        return Endpoint(service, self.full_name(), mapping.external_port)

    def service_info(self, name):
        service = self.get_service(name)
        if not service:
            raise Exception("service not found: {0}".format(name))
        return self.service_to_endpoint(service)

    def endpoints(self):
        return [self.service_to_endpoint(service) for service in self.service_config.load()]

    def sync(self):
        self.port_mapper.sync()

        services = self.service_config.load()
        domain = self.domain_config.load()
        if not domain:
            raise Exception("No token saved, need to call set_dns or get_dns_token first")

        services_data = []
        for service in services:
            mapping = self.port_mapper.get(service.port)
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

        data = {'token': domain.update_token, 'local_ip': self.local_ip, 'services': services_data}

        external_ip = self.port_mapper.external_ip()
        if not external_ip:
            self.logger.warn("Unable to get external ip")
        else:
            if IP(external_ip).iptype() != 'PUBLIC':
                external_ip = None
                self.logger.warn("External ip is not public")

        if external_ip:
            data['ip'] = external_ip
        else:
            self.logger.warn("Will try server side client ip detection")

        url = urljoin(self.insider_config.get_redirect_api_url(), "/domain/update")

        self.logger.debug('url: ' + url)
        json = convertible.to_json(data)
        self.logger.debug('request: ' + json)
        response = requests.post(url, json)

        self.check_http_error(response)