class AppSetup:

    def __init__(self, app_name, app_paths, nginx):
        self.app_name = app_name
        self.app_paths = app_paths
        self.nginx = nginx

    def get_app_dir(self):
        return self.app_paths.get_app_dir()

    def get_app_data_dir(self, remove_existing=False):
        return self.app_paths.get_app_data_dir(remove_existing)

    def register_web(self, port):
        self.nginx.add_app(self.app_name, port)

    def unregister_web(self):
        self.nginx.remove_app(self.app_name)

