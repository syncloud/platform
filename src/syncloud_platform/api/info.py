from syncloud_platform.di.injector import get_injector
from syncloud_platform.insider.device_info import DeviceInfo


def domain():
    injector = get_injector()
    return DeviceInfo(injector.user_platform_config, injector.port_config).domain()


def url(app=None):
    injector = get_injector()
    return DeviceInfo(injector.user_platform_config, injector.port_config).url(app)
