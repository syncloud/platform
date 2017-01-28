from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.application.apppaths import AppPaths
from syncloud_platform.application.appsetup import AppSetup
from syncloud_platform.injector import get_injector, PLATFORM_CONFIG_DIR


def get_app_paths(app_name, config_path=None):
    config = get_injector().platform_config
    return AppPaths(app_name, config)


def get_app_setup(app_name):
    injector = get_injector()
    app_paths = AppPaths(app_name, injector.platform_config)
    app_setup = AppSetup(app_name, app_paths, injector.nginx, injector.hardware, injector.info,
                         injector.device, injector.user_platform_config, injector.systemctl)
    return app_setup