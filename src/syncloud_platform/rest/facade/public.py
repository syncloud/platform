from syncloudlib import logger

from syncloud_platform.control import power


class Public:

    def __init__(self, platform_config, user_platform_config, snap, hardware, redirect_service,
                 log_aggregator, network):
        self.hardware = hardware
        self.platform_config = platform_config
        self.log = logger.get_logger('rest.public')
        self.user_platform_config = user_platform_config
        self.snap = snap
        self.www_dir = self.platform_config.www_root_public()
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator
        self.network=network
        
    def restart(self):
        power.restart()

    def shutdown(self):
        power.shutdown()

    def install(self, app_id):
        self.snap.install(app_id)

    def remove(self, app_id):
        return self.snap.remove(app_id)

    def upgrade(self, app_id):
        self.snap.upgrade(app_id)

    def disk_activate(self, device):
        return self.hardware.activate_disk(device)

    def send_logs(self, include_support):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user_token, logs, include_support)

    def network_interfaces(self):
        return self.network.interfaces()
