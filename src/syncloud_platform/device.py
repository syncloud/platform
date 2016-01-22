import getpass
import uuid

from syncloud_app import logger

from syncloud_platform.insider.util import protocol_to_port
from syncloud_platform.tools.chown import chown


class Device:

    def __init__(self, platform_config, user_platform_config, redirect_service,
                 port_drill_factory, common, sam, platform_cron, ldap_auth, event_trigger, tls):
        self.tls = tls
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.port_drill_factory = port_drill_factory
        self.sam = sam
        self.auth = ldap_auth
        self.platform_cron = platform_cron
        self.event_trigger = event_trigger
        self.common = common
        self.logger = logger.get_logger('Device')

    def activate(self,
                 redirect_email, redirect_password, user_domain,
                 device_user, device_password,
                 api_url=None, domain=None):

        if not api_url:
            api_url = 'http://api.syncloud.it'

        if not domain:
            domain = 'syncloud.it'

        self.logger.info("activate {0}, {1}, {2}, {3}, {4}".format(
            redirect_email, user_domain, device_user, api_url, domain))

        self.sam.update()

        self.user_platform_config.update_redirect(domain, api_url)
        user = self.redirect_service.get_user(redirect_email, redirect_password)
        self.user_platform_config.set_user_update_token(user.update_token)

        response_data = self.redirect_service.acquire(redirect_email, redirect_password, user_domain)
        self.user_platform_config.update_domain(response_data.user_domain, response_data.update_token)

        self.platform_cron.remove()
        self.platform_cron.create()

        self.set_access('http', False)

        self.logger.info("activating ldap")
        self.auth.reset(device_user, device_password)
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        self.tls.generate_certificate()

        self.logger.info("activation completed")

    def set_access(self, protocol, external_access):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        drill = self.get_drill(external_access)
        new_web_local_port = protocol_to_port(protocol)
        old_web_local_port = protocol_to_port(self.user_platform_config.get_protocol())

        try:
            drill.remove(old_web_local_port)
        except Exception, e:
            self.logger.error('Unable to remove port {0}: {1}'.format(old_web_local_port, e.message))

        try:
            drill.sync_new_port(new_web_local_port)
        except Exception, e:
            self.logger.error('Unable to add new port {0}: {1}'.format(new_web_local_port, e.message))

        self.redirect_service.sync(drill, update_token)
        self.user_platform_config.update_device_access(external_access, protocol)
        self.event_trigger.trigger_app_event_domain(self.platform_config.apps_root())

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            return

        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        self.redirect_service.sync(drill, update_token)

        if not getpass.getuser() == self.platform_config.cron_user():
            chown(self.platform_config.cron_user(), self.platform_config.data_dir())

    def send_logs(self):
        user_token = self.user_platform_config.get_user_update_token()
        logs = self.common.get_logs()
        self.redirect_service.send_log(user_token, logs)

    def get_drill(self, external_access):
        return self.port_drill_factory.get_drill(external_access)
