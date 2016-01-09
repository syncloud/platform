from syncloud_platform.di.injector import Injector
from syncloud_platform.insider.device_info import DeviceInfo


def domain():
    injector = Injector()
    return DeviceInfo(injector.user_platform_config, injector.port_config).domain()


def url(app=None):
    injector = Injector()
    return DeviceInfo(injector.user_platform_config, injector.port_config).url()
