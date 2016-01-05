import uuid

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.auth.ldapauth import LdapAuth
from syncloud_platform.insider import facade
from syncloud_app import logger
from syncloud_platform.tools.events import trigger_app_event_domain
from syncloud_platform.sam.stub import SamStub
from syncloud_platform.tools.tls import Tls
from syncloud_platform.insider.cron import PlatformCron


class Activator:
    def __init__(self, insider=None):
        if insider is None:
            insider = facade.get_insider()
        self.insider = insider
        self.logger = logger.get_logger('Activator')
        self.auth = LdapAuth()
        self.sam = SamStub()
        self.redirect_service = RedirectService()
        self.platform_config = PlatformConfig()
        self.user_platform_config = PlatformUserConfig(self.platform_config.get_user_config())
        self.platform_cron = PlatformCron(self.platform_config)


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

        self.redirect_service.set_info(domain, api_url)
        user = self.redirect_service.get_user(redirect_email, redirect_password)
        self.redirect_service.redirect_config.set_user_update_token(user.update_token)

        self.insider.acquire_domain(redirect_email, redirect_password, user_domain)

        self.platform_cron.remove()
        self.platform_cron.create()

        try:
            self.insider.add_main_device_service('http', False)
            trigger_app_event_domain(self.platform_config.apps_root())
        except Exception, e:
            self.logger.warn('upnp is not available ' + e.message)

        self.logger.info("activating ldap")
        self.auth.reset(device_user, device_password)
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        Tls().generate_certificate()

        self.user_platform_config.set_activated(True)
        self.logger.info("activation completed")

    def user_domain(self):
        return self.insider.user_domain()
