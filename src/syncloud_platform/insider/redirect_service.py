from urlparse import urljoin

import convertible
import requests
from IPy import IP
from syncloudlib import logger

from syncloud_platform.insider import util
from syncloud_platform.board import id
from syncloud_platform.gaplib import linux


class RedirectService:

    def __init__(self, user_platform_config, versions):
        self.versions = versions
        self.user_platform_config = user_platform_config

        self.logger = logger.get_logger('RedirectService')

    def get_user(self, email, password):
        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/user/get")
        response = requests.get(url, params={'email': email, 'password': password})
        util.check_http_error(response)
        user = convertible.from_json(response.text).data
        return user

    def send_log(self, user_update_token, logs, include_support):

        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/user/log")
        response = requests.post(url, {
            'token': user_update_token,
            'data': logs,
            'include_support': include_support
        })
        util.check_http_error(response)
        user = convertible.from_json(response.text)

        return user

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
        
    def sync(self, external_ip, web_port, web_local_port, web_protocol, update_token, external_access):
        
        map_local_address = not external_access
        
        version = self.versions.platform_version()
        local_ip = linux.local_ip()
        data = {
            'token': update_token,
            'platform_version': version,
            'local_ip': local_ip,
            'map_local_address': map_local_address,
            'web_protocol': web_protocol,
            'web_port': web_port,
            'web_local_port': web_local_port
        }

        if not external_ip:
            self.logger.warn("No external ip")
        else:
            iptype=IP(external_ip).iptype()
            if iptype != 'PUBLIC':
                external_ip = None
                self.logger.warn("External ip is not public: {0}".format(iptype))

        if not map_local_address:
            if external_ip:
                data['ip'] = external_ip
            else:
                self.logger.warn("Will try server side client ip detection")

        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/domain/update")

        self.logger.debug('url: ' + url)
        json = convertible.to_json(data)
        self.logger.info('request: ' + json)
        response = requests.post(url, json)

        util.check_http_error(response)