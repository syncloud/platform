from syncloudlib import logger

from syncloud_platform.control import power


class Public:

    def __init__(self, user_platform_config, redirect_service, log_aggregator, network):
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator
        self.network = network
        
    def restart(self):
        power.restart()

    def shutdown(self):
        power.shutdown()

    def send_logs(self, include_support):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user_token, logs, include_support)

    def network_interfaces(self):
        return self.network.interfaces()
