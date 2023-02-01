from urllib.parse import urljoin

from syncloudlib.json import convertible
import requests
from syncloudlib import logger

from syncloud_platform.insider import util


class RedirectService:

    def __init__(self, user_platform_config):
        self.user_platform_config = user_platform_config

        self.logger = logger.get_logger('RedirectService')

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

