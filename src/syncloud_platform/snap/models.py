from syncloudlib.json.convertible import Field, List


class App:
    id = Field()
    name = Field()
    required = Field(default=False)
    ui = Field(default=False)
    url = Field(default='')
    icon = Field(default=None)
    description = Field(default='No description given yet')
    enabled = Field(default=True)


class Apps:
    apps = Field(field_type=List(App))


class AppVersions:
    app = Field(field_type=App)
    current_version = Field()
    installed_version = Field()
