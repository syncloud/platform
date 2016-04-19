from syncloud_platform.di.injector import get_injector


def init(app_id, owner):
    return get_injector().hardware.init_app_storage(app_id, owner)
