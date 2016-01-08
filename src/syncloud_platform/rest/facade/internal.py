from syncloud_platform.tools import id


class Internal:
    def __init__(self, platform_config, device):
        self.device = device
        self.platform_config = platform_config
        self.www_dir = self.platform_config.www_root()

    def identification(self):
        return id.id()

    def activate(self,
                 redirect_email, redirect_password, user_domain,
                 device_user, device_password,
                 api_url=None, domain=None):

        self.device.activate(
            redirect_email,
            redirect_password,
            user_domain,
            device_user,
            device_password,
            api_url,
            domain
        )
