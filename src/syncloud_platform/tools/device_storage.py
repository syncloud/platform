class DeviceStorage:

    def __init__(self, hardware):
        self.hardware = hardware

    def init(self, app_id, owner):
        return self.hardware.init_app_storage(app_id, owner)
