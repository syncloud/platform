#!/opt/app/platform/python/bin/python

from syncloud_platform.injector import get_injector

injector = get_injector()

injector.device.sync_all()
user_platform_config = injector.user_platform_config
if user_platform_config.is_https() and user_platform_config.get_external_access():
    injector.tls.generate_real_certificate()
