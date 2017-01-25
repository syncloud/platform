from syncloud_platform.injector import get_injector
 
injector = get_injector()
hardware = injector.hardware
hardware.check_external_disk()
