from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.application.apppaths import AppPaths
from syncloud_platform.application.appsetup import AppSetup
from syncloud_platform.injector import get_injector

def get_app_paths(app_name, config_path=None):
    config = get_injector(config_path).platform_config()
    return AppPaths(app_name, config)

def get_app_setup(app_name):
    app_paths = get_app_paths(app_name)
    injector = get_injector()
    app_setup = AppSetup(app_name, app_paths, injector.nginx, injector.hardware, injector.info, injector.device, injector.user_platform_config)
    return app_setup