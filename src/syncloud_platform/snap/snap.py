from syncloud_app import logger
import json
import requests_unixsocket
from syncloud_platform.sam.models import AppVersions, App

SOCKET = "http+unix://%2Fvar%2Frun%2Fsnapd.socket"


class Snap:

    def __init__(self, platform_config, info):
        self.info = info
        self.platform_config = platform_config
        self.logger = logger.get_logger('Snap')

    def update(self, release=None):
        return None
        
    def install(self, app_id):
        self.logger.info('snap install')
        session = requests_unixsocket.Session()
        response = session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={'action': 'install'})
        self.logger.info("install response: {0}".format(response.text))
        

    def upgrade(self, app_id):
        self.logger.info('snap upgrade')
        session = requests_unixsocket.Session()
        response = session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), data={'action': 'install'})
        self.logger.info("install response: {0}".format(response.text))
        
    def status(self):
        self.logger.info('snap changes')
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/changes?select=in-progress'.format(SOCKET))
        self.logger.info("changes response: {0}".format(response.text))
        snapd_response = json.loads(response.text)
       
        if (snapd_response['status']) != 'OK':
            raise Exception(snapd_response['result']['message'])
            
        return len(snapd_response['result']) > 0

    def remove(self, app_id):
        self.logger.info('snap remove')
        session = requests_unixsocket.Session()
        session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), data={'action': 'remove'})

    def list(self):
        self.logger.info('snap list')
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/find?name=*'.format(SOCKET))
        self.logger.info("find response: {0}".format(response.text))
        return self.parse_response(response.text, lambda app: True)

    def find_in_store(self, app_id):
        self.logger.info('snap list')
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/find?name={0}'.format(SOCKET, app_id))
        self.logger.info("find app: {0}, response: {1}".format(app_id, response.text))
        
        found_apps = self.parse_response(response.text, lambda app: True)
        if (len(found_apps) == 0):
            self.logger.warn("No app found")
            retuen None
            
        if (len(found_apps) > 1):
           self.logger.warn("More than one app found")
        
        return found_apps[0]

    def user_apps(self):
        self.logger.info('snap user apps')
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/find?name=*'.format(SOCKET))
        self.logger.info("find response: {0}".format(response.text))
        return self.parse_response(response.text, lambda app: app['type'] == 'app')

    def installed_user_apps(self):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info("snaps response: {0}".format(response.text))
        return self.parse_response(response.text, lambda app: app['type'] == 'app')

    def installed_all_apps(self):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info("snaps response: {0}".format(response.text))
        return self.parse_response(response.text, lambda app: True)

    def find_installed(self, app_id):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps/{1}'.format(SOCKET, app_id))
        self.logger.info("snap response: {0}".format(response.text))
        snap_response = json.loads(response.text)
        existing_app = self.to_app(snap_response['result'])
        return existing_app
        
    def get_app(self, app_id):
        existing_app = self.find_installed(app_id)
        store_app = self.find_in_store(app_id)
        if store_app:
            existing_app.current_version = store_app.current_version
        retuen existing_app

    def parse_response(self, response_json, result_filter):
        response = json.loads(response_json)
        return [self.to_app(app) for app in response['result'] if result_filter(app)]

    def to_app(self, app):
    
        newapp = App()
        newapp.id = app['name']
        newapp.name = app['summary']
        newapp.url = self.info.url(newapp.id)
        newapp.icon = "http://apps.syncloud.org/releases/{0}/images/{1}-128.png".format(app['channel'], newapp.id)
        
        app_version = AppVersions()
        app_version.installed_version = app['version']
        app_version.current_version = app['version']
        app_version.app = newapp
        
        return app_version
