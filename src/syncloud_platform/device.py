import getpass
import uuid

from syncloud_app import logger

from syncloud_platform.auth.ldapauth import LdapAuth
from syncloud_platform.insider.cron import PlatformCron
from syncloud_platform.insider.util import protocol_to_port
from syncloud_platform.tools.chown import chown
from syncloud_platform.tools.events import trigger_app_event_domain
from syncloud_platform.tools.tls import Tls


class Device:

    def __init__(self, platform_config, user_platform_config, redirect_service, port_drill_factory, sam):
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.redirect_service = redirect_service
        self.port_drill_factory = port_drill_factory

        self.sam = sam
        self.auth = LdapAuth(self.platform_config)
        self.platform_cron = PlatformCron(self.platform_config)

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

        try:
            self.set_access('http', False)
        except Exception, e:
            self.logger.warn('upnp is not available ' + e.message)

        self.logger.info("activating ldap")
        self.auth.reset(device_user, device_password)
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        Tls().generate_certificate()

        self.logger.info("activation completed")

    def set_access(self, protocol, external_access):
        drill = self.get_drill(external_access)
        self.redirect_service.remove_service("server", drill)
        self.redirect_service.add_service("server", protocol, "server", protocol_to_port(protocol), drill)

        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            raise Exception("No update token saved, device is not activated yet")

        self.redirect_service.sync(drill, update_token)
        self.user_platform_config.update_device_access(external_access, protocol)
        trigger_app_event_domain(self.platform_config.apps_root())

    def sync_all(self):
        update_token = self.user_platform_config.get_domain_update_token()
        if update_token is None:
            raise Exception("No update token saved, device is not activated yet")

        external_access = self.user_platform_config.get_external_access()
        drill = self.get_drill(external_access)
        self.redirect_service.sync(drill, update_token)

        if not getpass.getuser() == self.platform_config.cron_user():
            chown(self.platform_config.cron_user(), self.platform_config.data_dir())

    def get_drill(self, external_access):
        return self.port_drill_factory.get_drill(external_access)
