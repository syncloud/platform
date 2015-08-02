from os import listdir
from urlparse import urljoin
from subprocess import check_output
from os.path import isfile, join
from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.insider import util
import requests
import convertible
from syncloud_platform.insider.config import RedirectConfig
from syncloud_platform.tools.app import get_app_data_root


class RedirectService:

    def __init__(self, data_root=None):
        self.log_root = PlatformConfig().get_log_root()
        if not data_root:
            data_root = get_app_data_root('platform')
        self.redirect_config = RedirectConfig(data_root)

    def get_user(self, email, password):
        url = urljoin(self.redirect_config.get_api_url(), "/user/get")
        response = requests.get(url, params={'email': email, 'password': password})
        util.check_http_error(response)
        user = convertible.from_json(response.text).data
        return user

    def send_log(self):

        log_files = [join(self.log_root, f) for f in listdir(self.log_root) if isfile(join(self.log_root, f))]
        log_files.append('/var/log/sam.log')

        logs = '\n----------------------\n'.join(map(self.read_log, log_files))

        url = urljoin(self.redirect_config.get_api_url(), "/user/log")
        response = requests.post(url, {'token': self.redirect_config.get_user_update_token(), 'data': logs})
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

    def set_info(self, domain, api_url):
        return self.redirect_config.update(domain, api_url)