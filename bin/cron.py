#!/opt/app/platform/python/bin/python

from syncloud_platform.injector import Injector

injector = Injector()

injector.device.sync_all()

injector.tls.generate_real_certificate()