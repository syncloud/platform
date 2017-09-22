from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.injector import get_injector, PLATFORM_CONFIG_DIR


def get_app_paths(app_name, config_path=None):
    return get_injector().get_app_paths(app_name)


def get_app_setup(app_name):
    return get_injector().get_app_setup(app_name)
