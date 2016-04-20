from syncloud_platform.application.apppaths import AppPaths
from syncloud_platform.config.config import PlatformConfig

def app_paths(app_name):
    return AppPaths(app_name, PlatformConfig())