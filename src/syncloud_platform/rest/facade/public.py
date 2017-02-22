from syncloud_app import logger

from syncloud_platform.gaplib.linux import pgrep, run_detached
from syncloud_platform.rest.model.app import app_from_sam_app
from syncloud_platform.control import power
from syncloud_platform.sam.stub import SAM_BIN_SHORT


class Public:

    def __init__(self, platform_config, user_platform_config, device, device_info, sam, hardware, redirect_service, log_aggregator, certbot_generator, port_mapper_factory, network, platforn_user_config):
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
        self.port_mapper_factory = port_mapper_factory
        self.network=network
        self.resize_script = self.platform_config.get_boot_extend_script()
        self.platforn_user_config = platforn_user_config

    def domain(self):
        return self.device_info.domain()

    def restart(self):
        power.restart()

    def shutdown(self):
        power.shutdown()

    def installed_apps(self):
        apps = [app_from_sam_app(a) for a in self.sam.installed_user_apps()]
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
        return pgrep(SAM_BIN_SHORT)

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

    def send_logs(self):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.log_aggregator.get_logs()
        self.redirect_service.send_log(user_token, logs)

    def regenerate_certificate(self):
        self.certbot_generator.generate_certificate()

    def port_mapper(self):
        enabled = self.user_platform_config.get_upnp()
        mapper = self.port_mapper_factory.provide_mapper()
        available = mapper is not None
        if available:
            message = 'Your router has {0} enabled, public ip: {1}'.format(mapper.name(),  mapper.external_ip())
        else:
            message = 'Your router does not have port mapping feature enabled at the moment'
        manual_public_ip = self.platforn_user_config.get_public_ip()
        return dict(available=available, enabled=enabled, message=message, public_ip = manual_public_ip)

    def network_interfaces(self):
        return self.network.interfaces()
