class App:
    def __init__(self, app_id, name, url):
        self.id = app_id
        self.name = name
        self.url = url


def app_from_sam_app(sam_app):
        url = sam_app.app.id
        # TODO: pip-less sam should not prefix apps with 'syncloud-'
        if 'syncloud-' in sam_app.app.id:
            url = sam_app.app.id[9:]
        return App(
            sam_app.app.id,
            sam_app.app.name,
            '/' + url)
