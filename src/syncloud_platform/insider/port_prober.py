from urlparse import urljoin
import requests
from syncloud_app import logger


class PortProber:

    def __init__(self, redirect_api_url, update_token):
        self.redirect_api_url = redirect_api_url
        self.update_token = update_token
        self.logger = logger.get_logger('PortProber')

    def probe_port(self, port, protocol):
        self.logger.info('probing {0}'.format(port))
        url = urljoin(self.redirect_api_url, "/probe/port")
        try:
            response = requests.get(url, params={'token': self.update_token, 'port': port, 'protocol': protocol})
            self.logger.info('response status_code: {0}'.format(response.status_code))
            self.logger.info('response text: {0}'.format(response.text))
            return response.status_code == 200 and response.text == 'OK'
        except Exception, e:
            self.logger.info('{0} is not reachable, error: {1}'.format(port, e.message))
            return False
