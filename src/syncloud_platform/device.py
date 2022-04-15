from syncloud_platform.config.config import WEB_ACCESS_PORT, WEB_PROTOCOL
from syncloudlib import logger



class Device:

    def __init__(self, user_platform_config, redirect_service, event_trigger):
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.event_trigger = event_trigger
        self.logger = logger.get_logger('Device')

    def set_access(self, external_access, manual_public_ip, manual_access_port):
        self.logger.info('set_access: external_access={0}'.format(external_access))

        if self.user_platform_config.is_redirect_enabled():
            self.redirect_service.sync(manual_public_ip, manual_access_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
                                       self.user_platform_config.get_domain_update_token(), external_access)

        self.user_platform_config.update_device_access(external_access, manual_public_ip, manual_access_port)
        self.event_trigger.trigger_app_event_domain()

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        external_access = self.user_platform_config.get_external_access()
        public_ip = self.user_platform_config.get_public_ip()
        manual_access_port = self.user_platform_config.get_manual_access_port()

        self.redirect_service.sync(public_ip, manual_access_port, WEB_ACCESS_PORT, WEB_PROTOCOL,
                                   update_token, external_access)
