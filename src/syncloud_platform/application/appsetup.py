class AppSetup:

    def __init__(self, app_name, app_paths, nginx, storage, device_info):
        self.app_name = app_name
        self.app_paths = app_paths
        self.nginx = nginx
        self.storage = storage
        self.device_info = device_info

    def get_install_dir(self):
        return self.app_paths.get_install_dir()

    def get_data_dir(self, remove_existing=False):
        return self.app_paths.get_data_dir(remove_existing)

    def register_web(self, port):
        self.nginx.add_app(self.app_name, port)

    def unregister_web(self):
        self.nginx.remove_app(self.app_name)

    def init_storage(self, user):
        return self.storage.init_app_storage(self.app_name, user)

    def device_domain_name(self):
        return self.device_info.domain()

    def app_url(self):
        return self.device_info.url(self.app_name)



