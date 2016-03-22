from syncloud_platform.di.injector import get_injector


def add_port(local_port):
    injector = get_injector()
    injector.device.add_port(local_port)


def remove_port(local_port):
    injector = get_injector()
    injector.device.remove_port(local_port)
