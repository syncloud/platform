from syncloud_platform.device import get_device
from syncloud_platform.tools import id


class Internal:
    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.www_dir = self.platform_config.www_root()

    def identification(self):
        return id.id()

    def activate(self,
             redirect_email, redirect_password, user_domain,
             device_user, device_password,
             api_url=None, domain=None):

        device = get_device()

        device.activate(
            redirect_email,
            redirect_password,
            user_domain,
            device_user,
            device_password,
            api_url,
            domain
        )
