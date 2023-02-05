import sqlite3
from os.path import isfile

from syncloudlib import logger

USER_CONFIG_DB = '/var/snap/platform/current/platform.db'

TRUE = 'true'
FALSE = 'false'


class PlatformUserConfig:
    def __init__(self, config_db=USER_CONFIG_DB):
        self.config_db = config_db
        self.log = logger.get_logger('PlatformUserConfig')

    def init_config(self):
        if not isfile(self.config_db):
            self.init_user_config()

    def is_activated(self):
        result = self._get('platform.activated')
        return to_bool(result)

    def get_web_secret_key(self):
        return self._get('platform.web_secret_key', 'default')

    def init_user_config(self):
            
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute("create table config (key varchar primary key, value varchar)")
        conn.close()

    def is_redirect_enabled(self):
        result = self._get('platform.redirect_enabled')
        return to_bool(result)

    def get_domain(self):
        return self._get('platform.domain')

    def get_user_domain(self):
        return self._get('platform.user_domain')

    def get_custom_domain(self):
        return self._get('platform.custom_domain')

    def get_redirect_domain(self):
        return self._get('redirect.domain', 'syncloud.it')

    def _get(self, key, default_value=None):
        self.init_config()
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute('select value from config where key = ?', (key,))
        row = cursor.fetchone()
        conn.close()
        if row:
            return row[0]
        
        return default_value
 
 
def to_bool(db_value, default=False):
    if db_value is None:
        return default
    return db_value == TRUE


def from_bool(bool_value):
    return TRUE if bool_value else FALSE
