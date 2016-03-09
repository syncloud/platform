from syncloud_platform.di.injector import get_injector


def add_port(local_port):
    injector = get_injector()
    external_access = injector.user_platform_config.get_external_access()
    drill = injector.device.get_drill(external_access)
    drill.sync_new_port(local_port)
