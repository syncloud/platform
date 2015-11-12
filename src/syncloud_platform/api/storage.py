from syncloud_platform.tools.hardware import Hardware


def init(app_id, owner):
    hardware = Hardware()
    return hardware.init_app_storage(app_id, owner)
