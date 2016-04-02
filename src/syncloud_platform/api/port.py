from syncloud_platform.di.injector import get_injector


def add_port(local_port, protocol):
    injector = get_injector()
    injector.device.add_port(local_port, protocol)


def remove_port(local_port, protocol):
    injector = get_injector()
    injector.device.remove_port(local_port, protocol)
