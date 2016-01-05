import uuid

from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.insider.facade import get_insider
from syncloud_platform.insider.cron import PlatformCron
from syncloud_platform.auth.ldapauth import LdapAuth
from syncloud_platform.sam.stub import SamStub
from syncloud_platform.tools.tls import Tls


class Activator:
    def __init__(self, platform_config, user_platform_config, insider, redirect_service):
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.insider = insider
        self.redirect_service = redirect_service

        self.logger = logger.get_logger('Activator')
        self.sam = SamStub()
        self.auth = LdapAuth(self.platform_config)
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

        self.user_platform_config.update_redirect(domain, api_url)
        user = self.redirect_service.get_user(redirect_email, redirect_password)
        self.user_platform_config.set_user_update_token(user.update_token)

        self.insider.acquire_domain(redirect_email, redirect_password, user_domain)

        self.platform_cron.remove()
        self.platform_cron.create()

        try:
            self.insider.add_main_device_service('http', False)
        except Exception, e:
            self.logger.warn('upnp is not available ' + e.message)

        self.logger.info("activating ldap")
        self.auth.reset(device_user, device_password)
        self.platform_config.set_web_secret_key(unicode(uuid.uuid4().hex))

        Tls().generate_certificate()

        self.user_platform_config.set_activated(True)
        self.logger.info("activation completed")

def get_activator():
    platform_config = PlatformConfig()
    user_platform_config = PlatformUserConfig(platform_config.get_user_config())
    insider = get_insider()
    redirect_service = RedirectService(platform_config, user_platform_config)

    return Activator(platform_config, user_platform_config, insider, redirect_service)
