import filecmp
import os
import uuid

from subprocess import check_output

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.auth.ldapauth import LdapAuth
from syncloud_platform.insider import facade
from syncloud_app import logger
from syncloud_platform.sam.stub import SamStub


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

        if filecmp.cmp(self.platform_config.get_ssl_certificate_file(), self.platform_config.get_default_ssl_certificate_file()):
            check_output('openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout {0} -out {1} -subj {2} 2>&1'.format(
                self.platform_config.get_ssl_key_file(),
                self.platform_config.get_ssl_certificate_file(),
                "/C=US/ST=Syncloud/L=Syncloud/O=Syncloud/OU=Syncloud/CN={0}.{1}".format(user_domain, domain)
            ), shell=True)
        else:
            self.logger.info("root ca exists, skipping")

        self.sam.update()
        # self.sam.upgrade_all()

        self.redirect_service.set_info(domain, api_url)
        user = self.redirect_service.get_user(redirect_email, redirect_password)
        self.redirect_service.redirect_config.set_user_update_token(user.update_token)

        self.insider.acquire_domain(redirect_email, redirect_password, user_domain)

        try:
            self.insider.add_main_device_service()
        except Exception, e:
            self.logger.warn('upnp is not available ' + e.message)

        self.logger.info("activating ldap")
        self.auth.reset(device_user, device_password)
        PlatformConfig().set_web_secret_key(unicode(uuid.uuid4().hex))

        PlatformUserConfig().set_activated(True)
        self.logger.info("activation completed")

    def user_domain(self):
        return self.insider.user_domain()
