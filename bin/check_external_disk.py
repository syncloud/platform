#!/opt/app/platform/python/bin/python 
 
from syncloud_platform.di.injector import Injector 
 
injector = Injector()
hardware = injector.hardware
hardware.check_external_disk()
