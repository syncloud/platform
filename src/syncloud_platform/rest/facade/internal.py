from syncloud_platform.tools import id


class Internal:
    def __init__(self, platform_config, device, redirect_service, log_aggregator):
        self.device = device
        self.platform_config = platform_config
        self.www_dir = self.platform_config.www_root()
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator

    def identification(self):
        return id.id()

    def activate(self, redirect_email, redirect_password, user_domain, device_username, device_password, main_domain):
        self.device.activate(
            redirect_email, redirect_password, user_domain,
            device_username, device_password, main_domain)

    def send_logs(self, redirect_email, redirect_password, main_domain):
        user = self.device.prepare_redirect(redirect_email, redirect_password, main_domain)
        
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user.update_token, logs)
