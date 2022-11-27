import sqlite3
import os
import uuid
from configparser import ConfigParser
from os.path import isfile, join
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

    def update_redirect(self, domain):
        self._upsert([
            ('redirect.domain', domain)
        ])

    def get_redirect_domain(self):
        return self._get('redirect.domain', 'syncloud.it')

    def get_redirect_api_url(self):
        return self._get('redirect.api_url', "https://api.{}".format(self.get_redirect_domain()))

    def get_user_update_token(self):
        return self._get('redirect.user_update_token')

    def get_user_email(self):
        return self._get('redirect.user_email')

    def set_custom_domain(self, custom_domain):
        self._upsert([
            ('platform.custom_domain', custom_domain)
        ])

    def set_deactivated(self):
        self._upsert([
            ('platform.activated', FALSE)
        ])

    def is_activated(self):
        result = self._get('platform.activated')
        return to_bool(result)

    def get_custom_domain(self):
        return self._get('platform.custom_domain')

    def get_domain(self):
        return self._get('platform.domain')

    def get_user_domain(self):
        return self._get('platform.user_domain')

    def update_domain(self, user_domain, domain_update_token):
        self._upsert([
            ('platform.user_domain', user_domain),
            ('platform.domain_update_token', domain_update_token)
        ])

    def is_redirect_enabled(self):
        result = self._get('platform.redirect_enabled')
        return to_bool(result)
        
    def set_redirect_enabled(self, enabled):
        self._upsert([
            ('platform.redirect_enabled', from_bool(enabled))
        ])
  
    def get_dkim_key(self):
        return self._get('dkim_key')

    def set_dkim_key(self, value):
        self._upsert([
            ('dkim_key', value)
        ])

    def get_manual_access_port(self):
        port = self._get('platform.manual_access_port')
        try:
            return int(port)
        except Exception as e:
            return None

    def set_manual_access_port(self, port):
        self._upsert([
             ('platform.manual_access_port', port)
        ])

    def get_web_secret_key(self):
        return self._get('platform.web_secret_key', 'default')

    def set_web_secret_key(self, value):
        self._upsert([
            ('platform.web_secret_key', value)
        ])

    def init_user_config(self):
            
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute("create table config (key varchar primary key, value varchar)")
        conn.close()

    def _upsert(self, key_values):
        self.init_config()
        conn = sqlite3.connect(self.config_db)
        with conn:
            for key, value in key_values:
                if value is not None:
                    self.log.info('setting {0}={1}'.format(key, value))
                    conn.execute('INSERT OR REPLACE INTO config VALUES (?, ?)', (key, value))
        conn.close() 

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
