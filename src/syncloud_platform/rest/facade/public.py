from syncloudlib import logger

from syncloud_platform.gaplib.linux import pgrep, run_detached
from syncloud_platform.rest.model.app import app_from_snap_app
from syncloud_platform.control import power


class Public:

    def __init__(self, platform_config, user_platform_config, device, device_info, snap, hardware, redirect_service,
                 log_aggregator, certbot_generator, port_mapper_factory, network, port_config):
        self.port_config = port_config
        self.hardware = hardware
        self.platform_config = platform_config
        self.log = logger.get_logger('rest.public')
        self.user_platform_config = user_platform_config
        self.device = device
        self.device_info = device_info
        self.snap = snap
        self.www_dir = self.platform_config.www_root_public()
        self.redirect_service = redirect_service
        self.log_aggregator = log_aggregator
        self.certbot_generator = certbot_generator
        self.port_mapper_factory = port_mapper_factory
        self.network=network
        self.resize_script = self.platform_config.get_boot_extend_script()

    def domain(self):
        return self.device_info.domain()

    def device_url(self):
        return self.device_info.url()

    def restart(self):
        power.restart()

    def shutdown(self):
        power.shutdown()

    def installed_apps(self):
        apps = [app_from_snap_app(a) for a in self.snap.installed_user_apps()]
        return apps

    def get_app(self, app_id):
        return self.snap.get_app(app_id)

    def list_apps(self):
        return self.snap.list()

    def install(self, app_id):
        self.snap.install(app_id)

    def remove(self, app_id):
        return self.snap.remove(app_id)

    def upgrade(self, app_id, channel, force):
        self.snap.upgrade(app_id, channel, force)

    def available_apps(self):
        return [app_from_snap_app(a) for a in self.snap.user_apps() if a.app.enabled]

    def port_mappings(self):
        return self.port_config.load()
    
    def access(self):
    
        upnp_enabled = self.user_platform_config.get_upnp()
        mapper = self.port_mapper_factory.provide_mapper()
        upnp_available = mapper is not None
        external_access = self.user_platform_config.get_external_access()

        result = dict(external_access=external_access,
                      upnp_available=upnp_available,
                      upnp_enabled=upnp_enabled,
                      upnp_message='not used')

        manual_public_ip = self.user_platform_config.get_public_ip()
        if manual_public_ip is not None:
            result['public_ip'] = manual_public_ip

        return result

    def set_access(self, upnp_enabled, external_access, public_ip, certificate_port, access_port):
        self.device.set_access(upnp_enabled, external_access, public_ip, certificate_port, access_port)

    def disk_activate(self, device):
        return self.hardware.activate_disk(device)

    def installer_status(self):
        return self.snap.status()

    def boot_extend_status(self):
        return pgrep(self.resize_script)

    def boot_extend(self):
        run_detached(self.resize_script, self.platform_config.get_platform_log(), self.platform_config.get_ssh_port())

    def disk_deactivate(self):
        return self.hardware.deactivate_disk()

    def disks(self):
        return self.hardware.available_disks()

    def boot_disk(self):
        return self.hardware.root_partition()

    def send_logs(self, include_support):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user_token, logs, include_support)

    def regenerate_certificate(self):
        self.certbot_generator.generate_certificate()

    def network_interfaces(self):
        return self.network.interfaces()
