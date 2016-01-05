from urlparse import urljoin

import requests
import convertible
from IPy import IP
from syncloud_app import logger
from syncloud_platform.insider import util
from syncloud_platform.insider.config import Service
from syncloud_platform.tools import id


class Dns:

    def __init__(self, service_config, local_ip, redirect_config, user_platform_config):
        self.redirect_config = redirect_config
        self.local_ip = local_ip
        self.service_config = service_config
        self.logger = logger.get_logger('dns')
        self.user_platform_config = user_platform_config

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
        response_data = convertible.from_json(response.text)
        self.user_platform_config.set_user_domain(response_data.user_domain)
        self.user_platform_config.set_update_token(response_data.update_token)
        return response_data

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

    def sync(self, port_drill):
        port_drill.sync()

        services = self.service_config.load()
        if not self.user_platform_config.is_activated():
            self.logger.info('nothing to sync yet, no dns configuration')
            return

        update_token = self.user_platform_config.get_update_token()
        if not update_token:
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
            'token': update_token,
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
