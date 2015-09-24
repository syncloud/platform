from urlparse import urljoin
import requests
from syncloud_app import logger


class PortProber:

    def __init__(self, domain_update_token, redirect_api_url):
        self.domain_update_token = domain_update_token
        self.redirect_api_url = redirect_api_url
        self.logger = logger.get_logger('PortProber')

    def probe_port(self, port):
        self.logger.info('probing {0}'.format(port))
        url = urljoin(self.redirect_api_url, "/probe/port")
        try:
            response = requests.get(url, params={'token': self.domain_update_token, 'port': port})
            return response.status_code == 200 and response.text == 'OK'
        except Exception, e:
            self.logger.info('{0} is not reachable'.format(port))
            return False
