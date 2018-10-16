#!/opt/app/platform/python/bin/python

from syncloud_platform.injector import get_injector

injector = get_injector()

injector.device.sync_all()
user_platform_config = injector.user_platform_config

# TODO: Should we generate real certificate when IP is public (no external access needed)
if user_platform_config.get_external_access() or not user_platform_config.get_redirect_enabled():
    injector.tls.generate_real_certificate()
