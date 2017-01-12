from syncloud_app import logger
import json
import requests_unixsocket
from syncloud_platform.rest.service_exception import ServiceException
from syncloud_platform.sam.models import AppVersions, App

SOCKET="http://unixsocket.requests-unixsocket.github.io//%2Fvar%2Frun%2Fsnapd.socket"

class Snap:

    def __init__(self, platform_config, info):
        self.info = info
        self.platform_config = platform_config
        self.logger = logger.get_logger('Snap')

    def update(self, release=None):
        return None
        
    def install(self, app_id):
        session = requests_unixsocket.Session()
        session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), data={'action': 'install'})

    def upgrade(self, app_id):
        session = requests_unixsocket.Session()
        session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), data={'action': 'install'})

    def remove(self, app_id):
        session = requests_unixsocket.Session()
        session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), data={'action': 'remove'})

    def list(self):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info(response)
        return response

    def user_apps(self):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info(response)
        return response

    def installed_user_apps(self):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info(response)
        return response

    def installed_all_apps(self):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info(response)
        apps = parse_snaps_response(response)
        return apps

    def get_app(self, app_id):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps/{1}'.format(SOCKET, app_id))
        self.logger.info(response)
        return response


def parse_snaps_response(response_json):
    response = json.loads(response_json)
    return [to_app(app) for app in response['result']]

def to_app(app):
    app_version = AppVersions()
    app_version.installed_version = app['version']
    app_version.current_version = app['version']
    app_version.app = App()
    app_version.app.id = app['name']
    return app_version

