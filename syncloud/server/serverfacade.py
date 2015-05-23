from os.path import join
from syncloud.server.auth import Auth
from syncloud.tools.facade import Facade
from syncloud.server.model import Credentials
from syncloud.apache.facade import ApacheFacade
from syncloud.sam.manager import get_sam
from syncloud.remote.remoteaccess import RemoteAccess
from syncloud.insider import facade
from syncloud.app import logger


class ServerFacade:
    def __init__(self, sam, insider, remote_access, apache):
        self.sam = sam
        self.insider = insider
        self.remote_access = remote_access
        self.apache = apache
        self.tools = Facade()
        self.logger = logger.get_logger('ServerFacade')
        self.auth = Auth()

    def activate(self, release, domain, api_url, redirect_email, redirect_password, user_domain):
        self.activate_new(release, domain, api_url, redirect_email, redirect_password, user_domain, 'syncloud', 'syncloud')

    def activate_new(self, release, domain, api_url, email, password, user_domain, device_user, device_password):

        self.reconfigure()

        self.logger.info("activate {0}, {1}, {2}, {3}, {4}".format(release, domain, api_url, email, user_domain))
        self.sam.update(release)
        self.sam.upgrade_all()
        self.insider.set_redirect_info(domain, api_url)
        self.insider.acquire_domain(email, password, user_domain)

        full_domain = "{0}.{1}".format(user_domain, domain)
        apache_ports = self.apache.activate(full_domain)
        self.insider.add_service("server", "http", "server", apache_ports.http, None)

        self.logger.info("reconfiguring installed apps")
        self.sam.reconfigure_installed_apps()

        self.logger.info("activating ldap")
        self.auth.reset(full_domain, device_user, device_password)

        credentials = _get_credentials(self.remote_access.enable())
        self.logger.info("activation completed")
        return credentials

    def reconfigure(self):
        http_conf = join(self.tools.usr_local_dir(), 'syncloud-server', 'apache', 'syncloud-server-http.conf')
        self.apache.add_http_site('server', http_conf)
        https_conf = join(self.tools.usr_local_dir(), 'syncloud-server', 'apache', 'syncloud-server-https.conf')
        self.apache.add_https_site('server', https_conf)
        self.apache.restart()

    def get_access(self):
        return _get_credentials(self.remote_access.add_certificate())

    def user_domain(self):
        return self.insider.user_domain()


def _get_credentials(private_key):
    return Credentials('root', 'syncloud', private_key)


def get_server(insider=None):
    sam = get_sam()
    if insider is None:
        insider = facade.get_insider()
    remote_access = RemoteAccess(insider)
    apache = ApacheFacade()
    return ServerFacade(sam, insider, remote_access, apache)