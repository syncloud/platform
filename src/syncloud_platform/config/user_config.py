import sqlite3
from ConfigParser import ConfigParser
from os.path import isfile, join
from syncloud_app import logger
from syncloud_platform.config import config

USER_CONFIG_FILE_OLD = join(config.DATA_DIR, 'user_platform.cfg')
USER_CONFIG_DB = join(config.DATA_DIR, 'platform.db')


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
        self.parser.read(self.filename)
        self.__set('redirect', 'user_update_token', user_update_token)
        self.__save()

    def get_user_update_token(self):
        self.parser.read(self.filename)
        return self.parser.get('redirect', 'user_update_token')

    def set_user_email(self, user_email):
        self.parser.read(self.filename)
        self.__set('redirect', 'user_email', user_email)
        self.__save()

    def get_user_email(self):
        self.parser.read(self.filename)
        return self.parser.get('redirect', 'user_email')

    def set_custom_domain(self, custom_domain):
        self.parser.read(self.filename)
        self.__set('platform', 'custom_domain', custom_domain)
        self.__save()

    def set_activated(self):
        self.parser.read(self.filename)
        self.__set('platform', 'activated', True)
        self.__save()

    def is_activated(self):

        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'activated'):
            return False
        return self.parser.getboolean('platform', 'activated')

    def get_custom_domain(self):
        self.parser.read(self.filename)
        if self.parser.has_option('platform', 'custom_domain'):
            return self.parser.get('platform', 'custom_domain')
        return None

    def get_user_domain(self):
        self.parser.read(self.filename)
        if self.parser.has_option('platform', 'user_domain'):
            return self.parser.get('platform', 'user_domain')
        return None

    def get_domain_update_token(self):
        self.parser.read(self.filename)
        if self.parser.has_option('platform', 'domain_update_token'):
            return self.parser.get('platform', 'domain_update_token')
        return None

    def update_domain(self, user_domain, domain_update_token):
        self.parser.read(self.filename)
        self.log.info('saving user_domain = {0}, domain_update_token = {0}'.format(user_domain, domain_update_token))
        self.__set('platform', 'user_domain', user_domain)
        self.__set('platform', 'domain_update_token', domain_update_token)
        self.__save()

    def get_external_access(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'external_access'):
            return False
        return self.parser.getboolean('platform', 'external_access')

    def is_redirect_enabled(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'redirect_enabled'):
            return True
        return self.parser.getboolean('platform', 'redirect_enabled')
    
    def set_redirect_enabled(self, enabled):
        self.parser.read(self.filename)
        self.__set('platform', 'redirect_enabled', enabled)
        self.__save()

    def update_device_access(self, upnp_enabled, external_access, public_ip, manual_certificate_port, manual_access_port):
        self.parser.read(self.filename)
        self.__set('platform', 'external_access', external_access)
        self.__set('platform', 'upnp', upnp_enabled)
        self.__set('platform', 'public_ip', public_ip)
        self.__set('platform', 'manual_certificate_port', manual_certificate_port)
        self.__set('platform', 'manual_access_port', manual_access_port)
        self.__save()

    def get_upnp(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'upnp'):
            return True
        return self.parser.getboolean('platform', 'upnp')

    def get_public_ip(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'public_ip'):
            return None
        return self.parser.get('platform', 'public_ip')

    def get_manual_certificate_port(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'manual_certificate_port'):
            return None
        return self.parser.get('platform', 'manual_certificate_port')

    def get_manual_access_port(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'manual_access_port'):
            return None
        return self.parser.get('platform', 'manual_access_port')

    def get_web_secret_key(self):
        self.parser.read(self.filename)
        if not self.parser.has_option('platform', 'web_secret_key'):
            return 'default'
        return self.parser.get('platform', 'web_secret_key')

    def set_web_secret_key(self, value):
        self.parser.read(self.filename)
        self.__set('platform', 'web_secret_key', value)
        self.__save()

    def __set(self, section, key, value):
        if not self.parser.has_section(section):
            self.parser.add_section(section)
        if value is None:
            self.parser.remove_option(section, key)
        else:
            self.parser.set(section, key, value)

    def init_user_config(self):
    
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute("create table config (key varchar primary key, value varchar)")
        conn.close()
    
    
    def migrate_user_config(self):
        if isfile(self.old_config_file):
            db.init_config_db()
            conn = sqlite3.connect(self.config_db)
            with conn:
                parser.read(self.old_config_file)
                # for section, key, valie in parser:
                #    _upsert(cirsor, '{0},{1}'.format(section, key), value)


    def _upsert(self,values):
        conn = sqlite3.connect(self.config_db)
        with conn:
            for key, value in values:
                self.log.info('setting {0}={1}'.format(key, value))
                conn.execute('INSERT OR REPLACE INTO config VALUES (?, ?)', (key, value))
        conn.close() 
     
 
    def _get(self, key, default_value):
        conn = sqlite3.connect(self.config_db)
        cursor = conn.cursor()
        cursor.execute('select value from config where key = ?', (key,))
        value, _ = cursor.fetchone()
        conn.close()
        if value:
            return value
        else:
            return default_value
 
    