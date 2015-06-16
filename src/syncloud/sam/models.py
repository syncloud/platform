from convertible import Field, List


class App:
    id = Field()
    name = Field()
    required = Field(default=False)
    ui = Field(default=False)


class Apps:
    apps = Field(field_type=List(App))


class AppVersions:
    def __init__(self, app, current_version, installed_version):
        self.app = app
        self.current_version = current_version
        self.installed_version = installed_version