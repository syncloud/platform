from syncloud_platform.di.injector import get_injector

def register_app(app, port):
    get_injector().nginx.add_app(app, port)


def unregister_app(app):
    get_injector().nginx.remove_app(app)