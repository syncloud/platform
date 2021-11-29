#!/snap/platform/current/python/bin/python
from syncloud_platform.gaplib import linux
from syncloud_platform.injector import get_injector

injector = get_injector()
user_platform_config = injector.user_platform_config

if not user_platform_config.is_activated():
    return

generate_real_certificate = True
if user_platform_config.is_redirect_enabled():
    injector.device.sync_all()
    if not user_platform_config.get_external_access():
        if not linux.is_ip_public(linux.local_ip()) and linux.local_ip_v6() is None:
            generate_real_certificate = False

if generate_real_certificate:
    injector.tls.generate_real_certificate()
