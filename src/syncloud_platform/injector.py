import logging
from miniupnpc import UPnP
from os import environ
from os.path import join

from syncloud_app import logger

from syncloud_platform.auth.ldapauth import LdapAuth
from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig, PLATFORM_APP_NAME, PLATFORM_CONFIG_DIR
from syncloud_platform.control.systemctl import Systemctl
from syncloud_platform.device import Device
from syncloud_platform.insider.cron import PlatformCron
from syncloud_platform.insider.device_info import DeviceInfo
from syncloud_platform.insider.natpmpc import NatPmpPortMapper
from syncloud_platform.insider.port_config import PortConfig
from syncloud_platform.insider.port_drill import PortDrillFactory
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.insider.upnpc import UpnpPortMapper
from syncloud_platform.log.aggregator import Aggregator
from syncloud_platform.rest.facade.internal import Internal
from syncloud_platform.rest.facade.public import Public
from syncloud_platform.sam.stub import SamStub
from syncloud_platform.snap.snap import Snap
from syncloud_platform.disks.lsblk import Lsblk
from syncloud_platform.disks.path_checker import PathChecker
from syncloud_platform.events import EventTrigger
from syncloud_platform.disks.hardware import Hardware
from syncloud_platform.control.nginx import Nginx
from syncloud_platform.certbot.certbot_generator import CertbotGenerator
from syncloud_platform.control.tls import Tls
from syncloud_platform.disks.udev import Udev
from syncloud_platform.application.apppaths import AppPaths
from syncloud_platform.versions import Versions

default_injector = None


def get_injector(config_dir=None):
    global default_injector
    if default_injector is None:
        default_injector = Injector(config_dir=config_dir)
    return default_injector


class Injector:
    def __init__(self, debug=False, config_dir=None):
        if not config_dir:
            config_dir = PLATFORM_CONFIG_DIR
        self.platform_config = PlatformConfig(config_dir=config_dir)

        if not logger.factory_instance:
            console = True if debug else False
            level = logging.DEBUG if debug else logging.INFO
            logger.init(level, console, join(self.platform_config.get_platform_log()))

        self.user_platform_config = PlatformUserConfig(self.platform_config.get_user_config())

        self.log_aggregator = Aggregator(self.platform_config)

        self.platform_app_paths = AppPaths(PLATFORM_APP_NAME, self.platform_config)
        self.platform_app_paths.get_data_dir()
        self.versions = Versions(self.platform_config)
        self.redirect_service = RedirectService(self.user_platform_config, self.versions)
        self.port_config = PortConfig(self.platform_app_paths.get_data_dir())

        self.nat_pmp_port_mapper = NatPmpPortMapper()
        self.upnp_port_mapper = UpnpPortMapper(UPnP())
        self.port_drill_factory = PortDrillFactory(self.user_platform_config, self.port_config,
                                                   self.nat_pmp_port_mapper, self.upnp_port_mapper)
        self.info = DeviceInfo(self.user_platform_config, self.port_config)
        if self.platform_config.get_installer() == 'sam':
            self.sam = SamStub(self.platform_config, self.info)
        else:
            self.sam = Snap(self.platform_config, self.info)
        self.platform_cron = PlatformCron(self.platform_config)
        self.systemctl = Systemctl(self.platform_config)
        self.ldap_auth = LdapAuth(self.platform_config, self.systemctl)
        self.event_trigger = EventTrigger(self.sam, self.platform_config)
        self.nginx = Nginx(self.platform_config, self.systemctl)
        self.certbot_genetator = CertbotGenerator(self.platform_config, self.user_platform_config, self.info, self.sam)
        self.tls = Tls(self.platform_config, self.user_platform_config, self.info, self.nginx, self.certbot_genetator)
        
        self.device = Device(self.platform_config, self.user_platform_config, self.redirect_service,
                             self.port_drill_factory, self.sam, self.platform_cron, self.ldap_auth,
                             self.event_trigger, self.tls)

        self.internal = Internal(self.platform_config, self.device, self.redirect_service, self.log_aggregator)
        self.path_checker = PathChecker(self.platform_config)
        self.lsblk = Lsblk(self.platform_config, self.path_checker)
        self.hardware = Hardware(self.platform_config, self.event_trigger,
                                 self.lsblk, self.path_checker, self.systemctl)

        self.public = Public(self.platform_config, self.user_platform_config, self.device, self.info, self.sam,
                             self.hardware, self.redirect_service, self.log_aggregator, self.certbot_genetator)
        self.udev = Udev(self.platform_config)
