import sqlite3
import os
from ConfigParser import ConfigParser
from os.path import isfile, join
from syncloud_app import logger
from syncloud_platform.config import config

USER_CONFIG_FILE_OLD = join(config.DATA_DIR, 'user_platform.cfg')
USER_CONFIG_DB = join(config.DATA_DIR, 'platform.db')

TRUE = 'true'
FALSE = 'false'

class PlatformUserConfig:
    def __init__(self, config_db=USER_CONFIG_DB, old_config_file=USER_CONFIG_FILE_OLD):
        self.config_db = config_db
        self.old_config_file = old_config_file
        self.log = logger.get_logger('PlatformUserConfig')

    def update_redirect(self, domain, api_url):
       self._upsert( [
            ('redirect.domain', domain),
            ('redirect.api_url', api_url)
        ])
    
    def get_redirect_domain(self):
        return self._get('redirect.domain', 'syncloud.it')

    def get_redirect_api_url(self):
        return self._get('redirect.api_url', 'http://api.syncloud.it')

    def set_user_update_token(self, user_update_token):
        self._upsert( [
            ('redirect.user_update_token', user_update_token)
        ])

    def get_user_update_token(self):
        return self._get('redirect.user_update_token')

    def set_user_email(self, user_email):
        self._upsert( [
            ('redirect.user_email', user_email)
        ])

    def get_user_email(self):
        return self._get('redirect.user_email')

    def set_custom_domain(self, custom_domain):
        self._upsert( [
            ('platform.custom_domain', custom_domain)
        ])

    def set_activated(self):
        self._upsert( [
            ('platform.activated', TRUE)
        ])

    def is_activated(self):
        result = self._get('platform.activated')
        return to_bool(result)

    def get_custom_domain(self):
        return self._get('platform.custom_domain')

    def get_user_domain(self):
        return self._get('platform.user_domain')

    def get_domain_update_token(self):
        return self._get('platform.domain_update_token')

    def update_domain(self, user_domain, domain_update_token):
        self._upsert( [
            ('platform.user_domain', user_domain),
            ('platform.domain_update_token', domain_update_token)
        ])

    def get_external_access(self):
        result = self._get('platform.external_access')
        return to_bool(result)

    def is_redirect_enabled(self):
        result = self._get('platform.redirect_enabled')
        return to_bool(result)
        
    def set_redirect_enabled(self, enabled):
        self._upsert( [
            ('platform.redirect_enabled', from_bool(enabled))
        ])

    def update_device_access(self, upnp_enabled, external_access, public_ip, manual_certificate_port, manual_access_port):
        self._upsert( [
            ('platform.external_access', from_bool(external_access)),
            ('platform.upnp', from_bool(upnp_enabled)),
            ('platform.public_ip', public_ip),
            ('platform.manual_certificate_port', manual_certificate_port),
            ('platform.manual_access_port', manual_access_port)
        ])

    def get_upnp(self):
        result = self._get('platform.upnp')
        return to_bool(result)

    def get_public_ip(self):
        return self._get('platform.public_ip')

    def get_manual_certificate_port(self):
        return self._get('platform.manual_certificate_port')

    def get_manual_access_port(self):
        return self._get('platform.manual_access_port')

    def get_web_secret_key(self):
        return self._get('platform.web_secret_key', 'default')

    def set_web_secret_key(self, value):
        self._upsert( [
            ('platform.web_secret_key', value)
        ])

    def init_user_config(self):
            
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute("create table config (key varchar primary key, value varchar)")
        conn.close()
    
    
    def migrate_user_config(self):
        if isfile(self.old_config_file):
            self.init_user_config()
            old_config = ConfigParser()
            old_config.read(self.old_config_file)
            for section in old_config.sections():
                for key, value in old_config.items(section):
                    db_value = from_bool(value == 'True') if value in ['True', 'False']  else value
                    self._upsert([
                        ('{0}.{1}'.format(section, key), db_value)
                    ])
            os.rename(self.old_config_file, self.old_config_file + '.bak')


    def _upsert(self,values):
        conn = sqlite3.connect(self.config_db)
        with conn:
            for key, value in values:
                self.log.info('setting {0}={1}'.format(key, value))
                conn.execute('INSERT OR REPLACE INTO config VALUES (?, ?)', (key, value))
        conn.close() 
     
 
    def _get(self, key, default_value=None):
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute('select value from config where key = ?', (key,))
        row = cursor.fetchone()
        conn.close()
        if row:
            return row[0]
        
        return default_value
 
 
def to_bool(db_value):
    if db_value is None:
        return False
    return db_value == TRUE

def from_bool(bool_value):
    return TRUE if bool_value else FALSE