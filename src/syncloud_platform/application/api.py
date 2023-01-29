from syncloud_platform.injector import get_injector


def get_app_setup(app_name):
    return get_injector().get_app_setup(app_name)
