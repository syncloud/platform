from os import listdir
from urlparse import urljoin
from subprocess import check_output
from os.path import isfile, join
import requests
import convertible
from IPy import IP

from syncloud_app import logger

from syncloud_platform.insider.config import Service
from syncloud_platform.insider import util
from syncloud_platform.tools import id

class RedirectService:

    def __init__(self, service_config, local_ip, user_platform_config, platform_config):
        self.local_ip = local_ip
        self.service_config = service_config
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.log_root = self.platform_config.get_log_root()

        self.logger = logger.get_logger('RedirectService')


    def get_user(self, email, password):
        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/user/get")
        response = requests.get(url, params={'email': email, 'password': password})
        util.check_http_error(response)
        user = convertible.from_json(response.text).data
        return user

    def send_log(self, user_update_token):

        log_files = [join(self.log_root, f) for f in listdir(self.log_root) if isfile(join(self.log_root, f))]
        log_files.append('/var/log/sam.log')

        logs = '\n----------------------\n'.join(map(self.read_log, log_files))

        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/user/log")
        response = requests.post(url, {'token': user_update_token, 'data': logs})
        util.check_http_error(response)
        user = convertible.from_json(response.text)

        return user

    def read_log(self, filename):
        log = 'file: {0}\n\n'.format(filename)
        if isfile(filename):
            log += check_output('tail -100 {0}'.format(filename), shell=True)
        else:
            log += '-- not found --'
        return log

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
        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/domain/acquire")
        response = requests.post(url, data)
        util.check_http_error(response)
        response_data = convertible.from_json(response.text)
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

    def sync(self, port_drill, update_token):
        port_drill.sync()
        services = self.service_config.load()

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

        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/domain/update")

        self.logger.debug('url: ' + url)
        json = convertible.to_json(data)
        self.logger.debug('request: ' + json)
        response = requests.post(url, json)

        util.check_http_error(response)
