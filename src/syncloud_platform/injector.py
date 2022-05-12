import logging
from os import environ
from os.path import join

from syncloud_platform.application.apppaths import AppPaths
from syncloud_platform.application.appsetup import AppSetup
from syncloud_platform.auth.ldapauth import LdapAuth
from syncloud_platform.config.config import PlatformConfig, PLATFORM_APP_NAME
from syncloud_platform.config.user_config import PlatformUserConfig
from syncloud_platform.control.nginx import Nginx
from syncloud_platform.control.systemctl import Systemctl
from syncloud_platform.disks.hardware import Hardware
from syncloud_platform.disks.lsblk import Lsblk
from syncloud_platform.disks.path_checker import PathChecker
from syncloud_platform.events import EventTrigger
from syncloud_platform.insider.device_info import DeviceInfo
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.log.aggregator import Aggregator
from syncloud_platform.network.network import Network
from syncloud_platform.rest.facade.public import Public
from syncloud_platform.snap.snap import Snap
from syncloudlib import logger

default_injector = None


def get_injector(debug=False):
    global default_injector
    if default_injector is None:
        config_dir = join(environ['SNAP'], 'config')
        default_injector = Injector(config_dir=config_dir, debug=debug)
    return default_injector


class Injector:
    def __init__(self, debug=False, config_dir=None):
        self.platform_config = PlatformConfig(config_dir=config_dir)

        if not logger.factory_instance:
            console = True if debug else False
            level = logging.DEBUG if debug else logging.INFO
            logger.init(level, console, join(self.platform_config.get_platform_log()))

        self.user_platform_config = PlatformUserConfig()

        self.log_aggregator = Aggregator(self.platform_config)

        self.platform_app_paths = AppPaths(PLATFORM_APP_NAME, self.platform_config)
        self.platform_app_paths.get_data_dir()
        self.redirect_service = RedirectService(self.user_platform_config)

        self.device_info = DeviceInfo(self.user_platform_config)
        self.snap = Snap(self.platform_config, self.device_info)
        self.systemctl = Systemctl(self.platform_config)
        self.ldap_auth = LdapAuth(self.platform_config, self.systemctl)
        self.event_trigger = EventTrigger(self.snap)
        self.nginx = Nginx(self.platform_config, self.device_info)

        self.path_checker = PathChecker(self.platform_config)
        self.lsblk = Lsblk(self.platform_config, self.path_checker)
        self.hardware = Hardware(self.platform_config, self.event_trigger,
                                 self.lsblk, self.path_checker, self.systemctl)
        self.network = Network()
        self.public = Public(self.platform_config, self.user_platform_config, self.device_info, self.snap,
                             self.hardware, self.redirect_service, self.log_aggregator,
                             self.network)

    def get_app_paths(self, app_name):
        return AppPaths(app_name, self.platform_config)
    
    def get_app_setup(self, app_name):
        return AppSetup(app_name, self.get_app_paths(app_name), self.nginx, self.hardware,
                        self.device_info, self.user_platform_config, self.systemctl)
