class App:
    def __init__(self, app_id, name, url):
        self.id = app_id
        self.name = name
        self.url = url


def app_from_sam_app(sam_app):
    return App(
        sam_app.app.id,
        sam_app.app.name,
        '/' + sam_app.app.id)
