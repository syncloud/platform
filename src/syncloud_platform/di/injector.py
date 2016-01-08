import logging

from os.path import join
from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig, PLATFORM_APP_NAME
from syncloud_platform.device import Device
from syncloud_platform.insider.port_config import PortConfig
from syncloud_platform.insider.port_drill import PortDrillFactory
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.insider.service_config import ServiceConfig
from syncloud_platform.rest.facade.common import Common
from syncloud_platform.rest.facade.internal import Internal
from syncloud_platform.rest.facade.public import Public
from syncloud_platform.sam.stub import SamStub
from syncloud_platform.tools import network
from syncloud_platform.tools.app import get_app_data_root


class Injector:
    def __init__(self, debug=False):
        self.platform_config = PlatformConfig()

        console = True if debug else False
        level = logging.DEBUG if debug else logging.INFO
        logger.init(level, console, join(self.platform_config.get_platform_log()))

        self.user_platform_config = PlatformUserConfig(self.platform_config.get_user_config())

        self.data_root = get_app_data_root(PLATFORM_APP_NAME)

        self.service_config = ServiceConfig(self.data_root)
        self.redirect_service = RedirectService(self.service_config, network.local_ip(), self.user_platform_config)
        self.port_config = PortConfig(self.data_root)

        self.port_drill_factory = PortDrillFactory(self.user_platform_config, self.port_config)
        self.sam = SamStub()
        self.device = Device(self.platform_config, self.user_platform_config, self.redirect_service, self.port_drill_factory, self.sam)

        self.common = Common(self.platform_config, self.user_platform_config, self.redirect_service)
        self.internal = Internal(self.platform_config, self.device)
        self.public = Public(self.platform_config, self.user_platform_config, self.device)
