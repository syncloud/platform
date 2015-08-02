import uuid
from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.server.auth import Auth
from syncloud_platform.tools.facade import Facade
from syncloud_platform.insider import facade
from syncloud_app import logger
from syncloud_platform.sam.stub import SamStub


class ServerFacade:
    def __init__(self, insider):
        self.insider = insider
        self.tools = Facade()
        self.logger = logger.get_logger('ServerFacade')
        self.auth = Auth()
        self.sam = SamStub()
        self.redirect_service = RedirectService()

    def activate(self,
                 redirect_email, redirect_password, user_domain,
                 device_user, device_password,
                 api_url=None, domain=None, release=None):

        if not api_url:
            api_url = 'http://api.syncloud.it'

        if not domain:
            domain = 'syncloud.it'

        self.logger.info("activate {0}, {1}, {2}, {3}, {4}, {5}".format(
            redirect_email, user_domain, device_user, release, api_url, domain))

        self.sam.update(release)
        # self.sam.upgrade_all()

        self.redirect_service.set_info(domain, api_url)
        user = self.redirect_service.get_user(redirect_email, redirect_password)
        self.redirect_service.redirect_config.set_user_update_token(user.update_token)

        self.insider.acquire_domain(redirect_email, redirect_password, user_domain)

        try:
            self.insider.add_service("server", "http", "server", 80, None)
        except Exception, e:
            self.logger.info('upnp is not available ' + e.message)

        self.logger.info("activating ldap")
        self.auth.reset(device_user, device_password)
        PlatformConfig().set_web_secret_key(unicode(uuid.uuid4().hex))

        PlatformUserConfig().set_activated(True)
        self.logger.info("activation completed")

    def user_domain(self):
        return self.insider.user_domain()


def get_server(insider=None):
    if insider is None:
        insider = facade.get_insider()
    return ServerFacade(insider)
