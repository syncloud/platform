from syncloudlib import logger

from syncloud_platform.rest.model.app import app_from_snap_app
from syncloud_platform.control import power


class Public:

    def __init__(self, platform_config, user_platform_config, device_info, snap, hardware, redirect_service,
                 log_aggregator, network):
        self.hardware = hardware
        self.platform_config = platform_config
        self.log = logger.get_logger('rest.public')
        self.user_platform_config = user_platform_config
        self.device_info = device_info
        self.snap = snap
        self.www_dir = self.platform_config.www_root_public()
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator
        self.network=network
        
    def domain(self):
        return self.device_info.domain()

    def device_url(self):
        return self.device_info.url()

    def restart(self):
        power.restart()

    def shutdown(self):
        power.shutdown()

    def get_app(self, app_id):
        return self.snap.get_app(app_id)

    def install(self, app_id):
        self.snap.install(app_id)

    def remove(self, app_id):
        return self.snap.remove(app_id)

    def upgrade(self, app_id):
        self.snap.upgrade(app_id)

    def disk_activate(self, device):
        return self.hardware.activate_disk(device)

    def installer_status(self):
        return self.snap.status()

    def disk_deactivate(self):
        return self.hardware.deactivate_disk()

    def send_logs(self, include_support):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user_token, logs, include_support)

    def network_interfaces(self):
        return self.network.interfaces()
