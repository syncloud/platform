from syncloud_platform.config.config import PLATFORM_APP_NAME
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.insider.service_config import ServiceConfig
from syncloud_platform.tools import network
from syncloud_platform.tools.app import get_app_data_root


html_prefix = '/server/html'
rest_prefix = '/server/rest'


class Common:
    def __init__(self, platform_config, user_platform_config):
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config

    def send_log(self):
        data_root = get_app_data_root(PLATFORM_APP_NAME)
        service_config = ServiceConfig(data_root)
        redirect_service = RedirectService(service_config, network.local_ip(), self.user_platform_config, self.platform_config)
        get_user_update_token = self.user_platform_config.get_user_update_token()
        redirect_service.send_log(get_user_update_token)
