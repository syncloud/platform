from syncloudlib import logger
import json
import requests_unixsocket
import requests
from syncloud_platform.snap.models import AppVersions, App
from syncloud_platform.gaplib.linux import pgrep, run_detached

SOCKET = "http+unix://%2Fvar%2Frun%2Fsnapd.socket"


class Snap:

    def __init__(self, platform_config, info):
        self.info = info
        self.platform_config = platform_config
        self.logger = logger.get_logger('Snap')

    def update(self, release=None):
        self.logger.info('snap update is not supported')
        return None

    def install(self, app_id):
        self.logger.info('snap install')
        session = requests_unixsocket.Session()
        response = session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={'action': 'install'})
        self.logger.info("install response: {0}".format(response.text))

    def upgrade(self, app_id, channel, force):
        self.logger.info('snap upgrade')
        if app_id == 'sam':
            self.upgrade_snapd(channel)
            return

        session = requests_unixsocket.Session()
        response = session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={
            'action': 'refresh',
            'channel': channel,
            'ignore-validation': force
        })
        self.logger.info("refresh response: {0}".format(response.text))
        snapd_response = json.loads(response.text)
        if (snapd_response['status']) != 'Accepted':
            raise Exception(snapd_response['result']['message'])

    def upgrade_snapd(self, channel):
        script = self.platform_config.get_snapd_upgrade_script()
     
        run_detached('{0} {1}'.format(script, channel),
                     self.platform_config.get_platform_log(),
                     self.platform_config.get_ssh_port())

    def snap_upgrade_status(self):
        return pgrep(self.platform_config.get_snapd_upgrade_script())

    def status(self):
        self.logger.info('snap changes')
        
        if self.snap_upgrade_status():
            self.logger.info("snapd upgrade is in progress")
            return True
            
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
        session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={'action': 'remove'})

    def list(self):
        installed_apps = self.installed_all_apps()
        store_apps = self.store_all_apps()
        apps = join_apps(installed_apps, store_apps)
        apps.append(self._installer())
        return apps

    def store_all_apps(self):
        self.logger.info('snap list')
        return [self._available_app(app) for app in self._available_snaps()]

    def find_in_store(self, app_id):
        self.logger.info('snap list')
        found_apps = [self._available_app(app) for app in self._available_snaps(app_id)]

        if len(found_apps) == 0:
            self.logger.warn("No app found")
            return None

        if len(found_apps) > 1:
            self.logger.warn("More than one app found")

        return found_apps[0]

    def user_apps(self):
        return [self._available_app(app) for app in self._available_snaps() if app['type'] == 'app']

    def _available_snaps(self, query='*'):
        self.logger.info('available snaps, query: {0}'.format(query))
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/find?name={1}'.format(SOCKET, query))
        self.logger.info("find response: {0}".format(response.text))
        snapd_response = json.loads(response.text)
        if query != "*" and snapd_response['status'] != 'OK':
            return []
        apps = [app for app in snapd_response['result'] if app['name'] != 'sam']
        return sorted(apps, key=lambda app: app['name'])

    def installed_user_apps(self):
        return [self._installed_app(app) for app in self._installed_snaps() if app['type'] == 'app']

    def installed_all_apps(self):
        return [self._installed_app(app) for app in self._installed_snaps()]

    def _installed_snaps(self):
        self.logger.info('installed snaps')
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps'.format(SOCKET))
        self.logger.info("snaps response: {0}".format(response.text))
        snap_response = json.loads(response.text)

        apps = snap_response['result']
        return sorted(apps, key=lambda app: app['name'])


    def _installer(self):
        channel = self.platform_config.get_channel()
        self.logger.info('system info')
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/system-info'.format(SOCKET))
        self.logger.info("system info response: {0}".format(response.text))
        snap_response = json.loads(response.text)

        version_response = requests.get('http://apps.syncloud.org/releases/{0}/snapd.version'.format(channel))

        return self.to_app(
            'sam',
            'Installer',
            channel,
            snap_response['result']['version'],
            version_response.text
        )

    def _installed_app(self, installed_app):
        return self.to_app(
            installed_app['name'],
            installed_app['summary'],
            installed_app['channel'],
            installed_app['version'],
            None)

    def _available_app(self, available_app):
        return self.to_app(
            available_app['name'],
            available_app['summary'],
            available_app['channel'],
            None,
            available_app['version'])

    def find_installed(self, app_id):
        session = requests_unixsocket.Session()
        response = session.get('{0}/v2/snaps/{1}'.format(SOCKET, app_id))
        self.logger.info("snap response: {0}".format(response.text))
        snap_response = json.loads(response.text)
        if snap_response['status-code'] == 404:
            return None
        app = snap_response['result']
        existing_app = self._installed_app(app)
        return existing_app

    def get_app(self, app_id):
        existing_app = self.find_installed(app_id)
        store_app = self.find_in_store(app_id)
        if not existing_app and not store_app:
            raise Exception("not found")

        if not store_app:
            return existing_app

        if not existing_app:
            return store_app

        existing_app.current_version = store_app.current_version
        return existing_app

    def to_app(self, id, name, channel, installed_version, store_version):

        new_app = App()
        new_app.id = id
        new_app.name = name
        new_app.url = self.info.url(id)
        new_app.icon = "/rest/app_image?channel={0}&app={1}".format(channel, id)

        app_version = AppVersions()
        app_version.installed_version = installed_version
        app_version.current_version = store_version
        app_version.app = new_app

        return app_version


def join_apps(installed_apps, store_apps):
    all_apps = dict([(app.app.id, app) for app in installed_apps])
    for store_app in store_apps:
        if store_app.app.id in all_apps:
            all_apps[store_app.app.id].current_version = store_app.current_version
        else:
            all_apps[store_app.app.id] = store_app

    return all_apps.values()
