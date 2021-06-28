from syncloud_platform.config.config import WEB_CERTIFICATE_PORT, WEB_ACCESS_PORT, WEB_PROTOCOL
from syncloudlib import logger

http_network_protocol = 'TCP'


class Device:

    def __init__(self, platform_config, user_platform_config, redirect_service,
                 port_drill_factory, event_trigger):
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.port_drill_factory = port_drill_factory
        self.event_trigger = event_trigger
        self.logger = logger.get_logger('Device')

    def set_access(self, upnp_enabled, external_access, manual_public_ip, manual_certificate_port, manual_access_port):
        self.logger.info('set_access: external_access={0}'.format(external_access))

        drill = self.port_drill_factory.get_drill(upnp_enabled, external_access, manual_public_ip,
                                                  manual_certificate_port, manual_access_port)

        if drill is None:
            self.logger.error('Will not change access mode. Was not able to get working port mapper.')
            return

        drill.sync_new_port(WEB_CERTIFICATE_PORT, http_network_protocol)
        mapping = drill.sync_new_port(WEB_ACCESS_PORT, http_network_protocol)
        router_port = None
        if mapping:
            router_port = mapping.external_port

        external_ip = drill.external_ip()

        if self.user_platform_config.is_redirect_enabled():
            self.redirect_service.sync(external_ip, router_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
                                       self.user_platform_config.get_domain_update_token(), external_access)

        self.user_platform_config.update_device_access(upnp_enabled, external_access,
                                                       manual_public_ip, manual_certificate_port, manual_access_port)
        self.event_trigger.trigger_app_event_domain()

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_certificate_port = self.user_platform_config.get_manual_certificate_port()
        manual_access_port = self.user_platform_config.get_manual_access_port()
        port_drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip,
                                                       manual_certificate_port, manual_access_port)
        try:
            port_drill.sync_existing_ports()
        except Exception as e:
            self.logger.error('Unable to sync port mappings: {0}'.format(str(e)))

        router_port = None
        mapping = port_drill.get(WEB_ACCESS_PORT, http_network_protocol)
        if mapping:
            router_port = mapping.external_port
        
        external_ip = port_drill.external_ip()
        
        self.redirect_service.sync(external_ip, router_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
                                   update_token, external_access)

    def add_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_certificate_port = self.user_platform_config.get_manual_certificate_port()
        manual_access_port = self.user_platform_config.get_manual_access_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip,
                                                  manual_certificate_port, manual_access_port)
        drill.sync_new_port(local_port, protocol)

    def remove_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        upnp = self.user_platform_config.get_upnp()
        public_ip = self.user_platform_config.get_public_ip()
        manual_certificate_port = self.user_platform_config.get_manual_certificate_port()
        manual_access_port = self.user_platform_config.get_manual_access_port()
        drill = self.port_drill_factory.get_drill(upnp, external_access, public_ip,
                                                  manual_certificate_port, manual_access_port)
        drill.remove(local_port, protocol)

