from os import listdir
from urlparse import urljoin
from subprocess import check_output
from os.path import isfile, join
from syncloud_platform.insider import util
import requests
import convertible


class RedirectService:

    def __init__(self, redirect_config, log_root):
        self.log_root = log_root
        self.redirect_config = redirect_config

    def get_user(self, email, password):
        url = urljoin(self.redirect_config.get_api_url(), "/user/get")
        response = requests.get(url, params={'email': email, 'password': password})
        util.check_http_error(response)
        user = convertible.from_json(response.text).data
        return user

    def send_log(self):

        log_files = [join(self.log_root, f) for f in listdir(self.log_root) if isfile(join(self.log_root, f))]
        logs = '\n----------------------\n'.join(map(self.read_log, log_files))

        url = urljoin(self.redirect_config.get_api_url(), "/user/log")
        response = requests.post(url, {'token': self.redirect_config.get_user_update_token(), 'data': logs})
        util.check_http_error(response)
        user = convertible.from_json(response.text)

        return user

    def read_log(self, filename):
        log = 'file: {0}\n'.format(filename)
        if isfile(filename):
            log += check_output('tail -100 {0}'.format(filename), shell=True)
        else:
            log += '-- not found --'
        return log

