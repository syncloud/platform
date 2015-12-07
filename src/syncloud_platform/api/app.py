from syncloud_platform.tools.nginx import Nginx


def register_app(app, port):
    Nginx().add_app(app, port)


def unregister_app(app):
    Nginx().remove_app(app)