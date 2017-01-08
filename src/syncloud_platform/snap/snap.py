from syncloud_app import logger
import jsonpickle
import requests
from syncloud_platform.rest.service_exception import ServiceException


class Snap:

    def __init__(self, platform_config, info):
        self.info = info
        self.platform_config = platform_config
        self.logger = logger.get_logger('Snap')

    def update(self, release=None):
        return None
        
    def install(self, app_id):
        requests.post('https://localhost:8181/v2/snaps/{0}'.format(app_id), data={'action': 'install'})

    def upgrade(self, app_id):
        requests.post('https://localhost:8181/v2/snaps/{0}'.format(app_id), data={'action': 'install'})

    def remove(self, app_id):
        requests.post('https://localhost:8181/v2/snaps/{0}'.format(app_id), data={'action': 'remove'})

    def list(self):
        return requests.get('https://localhost:8181/v2/snaps')

    def user_apps(self):
        return requests.get('https://localhost:8181/v2/snaps')

    def installed_user_apps(self):
        return requests.get('https://localhost:8181/v2/snaps')

    def installed_all_apps(self):
        return requests.get('https://localhost:8181/v2/snaps')

    def get_app(self, app_id):
        requests.get('https://localhost:8181/v2/snaps/{0}'.format(app_id))
