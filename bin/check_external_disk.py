#!/opt/app/platform/python/bin/python 
 
from syncloud_platform.di.injector import get_injector 
 
injector = get_injector()
hardware = injector.hardware
hardware.check_external_disk()
