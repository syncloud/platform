from syncloud_platform.di.injector import Injector


def init(app_id, owner):
    return Injector().storage.init(app_id, owner)
