from syncloud_platform.insider.config import Service
from syncloud_platform.insider.service_config import ServiceConfig
from syncloud_platform.tools.nginx import Nginx


def register_app(app, port):
    Nginx().add_app(app, port)
    ServiceConfig().add_or_update(Service(app, "http", "type", 0, app))


def unregister_app(app):
    Nginx().remove_app(app)
    ServiceConfig().remove(app)
