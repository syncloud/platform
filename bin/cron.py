#!/opt/app/platform/python/bin/python

from syncloud_platform.injector import Injector

injector = Injector()

injector.device.sync_all()
user_platform_config = injector.user_platform_config
if user_platform_config.get_protocol() == 'https' and user_platform_config.get_external_access():
    injector.tls.generate_real_certificate()
