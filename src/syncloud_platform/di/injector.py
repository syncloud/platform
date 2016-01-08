from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.rest.facade.common import Common
from syncloud_platform.rest.facade.internal import Internal
from syncloud_platform.rest.facade.public import Public


class Injector:
    def __init__(self):
        self.platform_config = PlatformConfig()
        logger.init(filename=self.platform_config.get_platform_log())
        self.user_platform_config = PlatformUserConfig(self.platform_config.get_user_config())
        self.internal = Internal(self.platform_config)
        self.public = Public(self.platform_config, self.user_platform_config)
        self.common = Common(self.platform_config, self.user_platform_config)
