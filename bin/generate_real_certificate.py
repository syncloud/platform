#!/snao/platform/current/python/bin/python
from syncloud_platform.injector import get_injector

injector = get_injector()
injector.tls.generate_real_certificate()
