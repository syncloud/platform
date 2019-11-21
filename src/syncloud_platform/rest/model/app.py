class App:
    def __init__(self, app_id, name, icon, url):
        self.id = app_id
        self.name = name
        self.url = url
        self.icon = icon


def app_from_snap_app(snap_app):
    return App(
        snap_app.app.id,
        snap_app.app.name,
        snap_app.app.icon,
        snap_app.app.url)
