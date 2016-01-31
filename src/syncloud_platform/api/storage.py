from syncloud_platform.di.injector import get_injector


def init(app_id, owner):
    return get_injector().storage.init(app_id, owner)
