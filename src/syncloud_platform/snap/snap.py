from syncloudlib import logger
import json
import requests_unixsocket

SOCKET = "http+unix://%2Fvar%2Frun%2Fsnapd.socket"


class Snap:

    def __init__(self):
        self.logger = logger.get_logger('Snap')

    def install(self, app_id):
        self.logger.info('snap install')
        session = requests_unixsocket.Session()
        response = session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={'action': 'install'})
        self.logger.info("install response: {0}".format(response.text))

    def upgrade(self, app_id):
        self.logger.info('snap upgrade')
        session = requests_unixsocket.Session()
        response = session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={
            'action': 'refresh'
        })
        self.logger.info("refresh response: {0}".format(response.text))
        snapd_response = json.loads(response.text)
        if (snapd_response['status']) != 'Accepted':
            raise Exception(snapd_response['result']['message'])

    def remove(self, app_id):
        self.logger.info('snap remove')
        session = requests_unixsocket.Session()
        session.post('{0}/v2/snaps/{1}'.format(SOCKET, app_id), json={'action': 'remove'})
