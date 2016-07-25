from syncloud_app import logger

from syncloud_platform.rest.props import html_prefix
from syncloud_platform.rest.model.app import app_from_sam_app, App

from syncloud_platform.control import power
from syncloud_platform.certbot import certbot_generator
class Public:

    def __init__(self, platform_config, user_platform_config, device, device_info, sam, hardware, redirect_service, log_aggregator, certbot_generator):
        self.hardware = hardware
        self.platform_config = platform_config
        self.log = logger.get_logger('rest.public')
        self.user_platform_config = user_platform_config
        self.device = device
        self.device_info = device_info
        self.sam = sam
        self.www_dir = self.platform_config.www_root_public()
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator
        self.certbot_generator = certbot_generator

    def domain(self):
        return self.device_info.domain()

    def restart(self):
        power.restart()

    def shutdown(self):
        power.shutdown()

    def installed_apps(self):
        apps = [app_from_sam_app(a) for a in self.sam.installed_user_apps()]

        # TODO: Hack to add system apps, need to think about it
        apps.append(App('store', 'App Store', None, html_prefix + '/store.html'))
        apps.append(App('settings', 'Settings', None, html_prefix + '/settings.html'))
        return apps

    def get_app(self, app_id):
        return self.sam.get_app(app_id)

    def list_apps(self):
        return self.sam.list()

    def install(self, app_id):
        self.sam.install(app_id)

    def remove(self, app_id):
        return self.sam.remove(app_id)

    def upgrade(self, app_id):
        self.sam.upgrade(app_id)

    def update(self):
        return self.sam.update()

    def available_apps(self):
        return [app_from_sam_app(a) for a in self.sam.user_apps()]

    def access(self):
        external_access = self.user_platform_config.get_external_access()
        protocol = self.user_platform_config.get_protocol()
        return dict(external_access=external_access, protocol=protocol)

    def external_access(self):
        return self.user_platform_config.get_external_access()

    def external_access_enable(self, external_access):
        self.device.set_access(self.user_platform_config.get_protocol(), external_access)

    def protocol(self):
        return self.user_platform_config.get_protocol()

    def set_protocol(self, protocol):
        self.device.set_access(protocol, self.user_platform_config.get_external_access())

    def disk_activate(self, device):
        return self.hardware.activate_disk(device)

    def system_upgrade(self):
        self.sam.upgrade('platform')

    def sam_upgrade(self):
        self.sam.upgrade('sam')

    def sam_status(self):
        return self.sam.is_running()

    def disk_deactivate(self):
        return self.hardware.deactivate_disk()

    def disks(self):
        return self.hardware.available_disks()

    def send_logs(self):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user_token, logs)

    def regenerate_certificate(self):
        self.certbot_generator.generate_certificate()
