from urlparse import urljoin
import requests
from syncloudlib import logger
import json
from IPy import IP


class PortProber:

    def __init__(self, redirect_api_url, update_token):
        self.redirect_api_url = redirect_api_url
        self.update_token = update_token
        self.logger = logger.get_logger('PortProber')

    def probe_port(self, port, protocol, ip):
        self.logger.info('probing {0}'.format(port))
        url = urljoin(self.redirect_api_url, "/probe/port_v2")
        try:
            request = {'token': self.update_token, 'port': port, 'protocol': protocol}
            if ip:
                iptype=IP(ip).iptype()
                if iptype != 'PUBLIC':
                    return False, 'IP: {0} is not public'.format(ip)
                request['ip'] = ip

            response = requests.get(url, params=request)
            self.logger.info('response status_code: {0}'.format(response.status_code))
            self.logger.info('response text: {0}'.format(response.text))
            result = json.loads(response.text)
            if response.status_code == 200 and result['message'] == 'OK':
                return True, ''
            else:
                external_device_ip = result['device_ip']
                ip_version = IP(external_device_ip).version()
                return False, 'using device public IP: "{0}" which is IPv{1}'.format(external_device_ip, ip_version)
                    
        except Exception, e:
            self.logger.info('{0} is not reachable, error: {1}'.format(port, e.message))
            return False, 'unable to validate external port: {0}'.format(e.message)


class NoneProber:
    def probe_port(self, port, protocol, ip):
        return True, ''