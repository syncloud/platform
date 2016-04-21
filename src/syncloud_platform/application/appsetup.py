class AppSetup:

    def __init__(self, app_name, app_paths, nginx, storage):
        self.app_name = app_name
        self.app_paths = app_paths
        self.nginx = nginx
        self.storage = storage

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

