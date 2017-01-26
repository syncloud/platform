import getpass
import uuid

from syncloud_app import logger

from syncloud_platform.insider.util import protocol_to_port
from syncloud_platform.gaplib import fs

http_network_protocol = 'TCP'


class Device:

    def __init__(self, platform_config, user_platform_config, redirect_service,
                 port_drill_factory, sam, platform_cron, ldap_auth, event_trigger, tls):
        self.tls = tls
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.port_drill_factory = port_drill_factory
        self.sam = sam
        self.auth = ldap_auth
        self.platform_cron = platform_cron
        self.event_trigger = event_trigger
        self.logger = logger.get_logger('Device')

    def prepare_redirect(self, redirect_email, redirect_password, main_domain):

        redirect_api_url = 'http://api.' + main_domain

        self.logger.info("prepare redirect {0}, {1}".format(redirect_email, redirect_api_url))

        self.sam.update()
        self.user_platform_config.update_redirect(main_domain, redirect_api_url)
        self.user_platform_config.set_user_email(redirect_email)

        user = self.redirect_service.get_user(redirect_email, redirect_password)
        return user

    def activate(self, redirect_email, redirect_password, user_domain, device_username, device_password, main_domain):

        self.logger.info("activate {0}, {1}".format(user_domain, device_username))
           
        user = self.prepare_redirect(redirect_email, redirect_password, main_domain)
        self.user_platform_config.set_user_update_token(user.update_token)

        response_data = self.redirect_service.acquire(redirect_email, redirect_password, user_domain)
        self.user_platform_config.update_domain(response_data.user_domain, response_data.update_token)

        self.platform_cron.remove()
        self.platform_cron.create()

        self.set_access('http', False)

        self.logger.info("activating ldap")
        fix_permissions = self.platform_config.get_installer() == 'sam'
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        self.tls.generate_self_signed_certificate()

        self.auth.reset(device_username, device_password, fix_permissions)

        self.logger.info("activation completed")

    def set_access(self, protocol, external_access):
        self.logger.info('set_access: protocol={0}, external_access={1}'.format(protocol, external_access))

        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        new_web_local_port = protocol_to_port(protocol)
        old_web_local_port = protocol_to_port(self.user_platform_config.get_protocol())

        if protocol == 'https' and external_access:
            self.tls.generate_real_certificate()

        drill = self.get_drill(external_access)

        if drill is None:
            self.logger.error('Will not change access mode. Was not able to get working port drill.')
            return
        try:
            drill.remove(old_web_local_port, http_network_protocol)
        except Exception, e:
            self.logger.error('Unable to remove port {0}: {1}'.format(old_web_local_port, e.message))

        try:
            drill.sync_new_port(new_web_local_port, http_network_protocol)
        except Exception, e:
            self.logger.error('Unable to add new port {0}: {1}'.format(new_web_local_port, e.message))

        self.redirect_service.sync(drill, update_token, protocol, external_access, http_network_protocol)
        self.user_platform_config.update_device_access(external_access, protocol)
        self.event_trigger.trigger_app_event_domain()

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        web_protocol = self.user_platform_config.get_protocol()
        self.redirect_service.sync(drill, update_token, web_protocol, external_access, http_network_protocol)

        if not getpass.getuser() == self.platform_config.cron_user():
            fs.chownpath(self.platform_config.data_dir(), self.platform_config.cron_user())

    def add_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        drill.sync_new_port(local_port, protocol)

    def remove_port(self, local_port, protocol):
        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        drill.remove(local_port, protocol)

    def get_drill(self, external_access):
        return self.port_drill_factory.get_drill(external_access)
