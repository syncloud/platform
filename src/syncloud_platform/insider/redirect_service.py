from os import listdir
from urlparse import urljoin
from subprocess import check_output
from os.path import isfile, join
from syncloud_platform.insider import util
import requests
import convertible


class RedirectService:

    def __init__(self, platform_config, user_platform_config):
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.log_root = self.platform_config.get_log_root()

    def get_user(self, email, password):
        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/user/get")
        response = requests.get(url, params={'email': email, 'password': password})
        util.check_http_error(response)
        user = convertible.from_json(response.text).data
        return user

    def send_log(self):

        log_files = [join(self.log_root, f) for f in listdir(self.log_root) if isfile(join(self.log_root, f))]
        log_files.append('/var/log/sam.log')

        logs = '\n----------------------\n'.join(map(self.read_log, log_files))

        url = urljoin(self.user_platform_config.get_redirect_api_url(), "/user/log")
        response = requests.post(url, {'token': self.user_platform_config.get_user_update_token(), 'data': logs})
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