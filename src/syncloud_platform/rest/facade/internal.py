class Internal:
    def __init__(self, platform_config, device, redirect_service, log_aggregator):
        self.device = device
        self.platform_config = platform_config
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator

    def send_logs(self, redirect_email, redirect_password, main_domain):
        user = self.device.prepare_redirect(redirect_email, redirect_password, main_domain)
        
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user.update_token, logs, True)
