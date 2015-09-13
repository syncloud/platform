from ConfigParser import ConfigParser
from os.path import dirname, join
import os
from syncloud_app import logger
import convertible
from syncloud_platform.config.config import PLATFORM_CONFIG_DIR


class Domain:

    def __init__(self, user_domain, update_token):
        self.update_token = update_token
        self.user_domain = user_domain


class Port:

    def __init__(self, local_port, external_port):
        self.local_port = local_port
        self.external_port = external_port


class Service:

    def __init__(self, name, protocol, type, port, url):
        self.name = name
        self.protocol = protocol
        self.type = type
        self.port = port
        self.url = url


class DomainConfig:

    def __init__(self, filename):
        self.filename = filename
        self.logger = logger.get_logger('insider.DomainConfig')

    def load(self):
        if not os.path.isfile(self.filename):
            raise Exception('{} does not exist'.format(self.filename))
        obj = convertible.read_json(self.filename)
        return obj

    def save(self, obj):
        self.logger.info('saving config={0}'.format(self.filename))
        convertible.write_json(self.filename, obj)

    def remove(self):
        if os.path.isfile(self.filename):
            os.remove(self.filename)

REDIRECT_CONFIG_NAME = 'redirect.cfg'


class RedirectConfig:
    def __init__(self, config_dir):
        self.parser = ConfigParser()
        self.filename = join(config_dir, REDIRECT_CONFIG_NAME)
        self.logger = logger.get_logger('insider.RedirectConfig')

    def update(self, domain, api_url):
        self.parser.read(self.filename)
        self.logger.info('log file: ' + self.filename)
        self.logger.info('setting domain={0}, api_url={1}'.format(domain, api_url))
        if not self.parser.has_section('redirect'):
            self.parser.add_section('redirect')

        self.parser.set('redirect', 'domain', domain)
        self.parser.set('redirect', 'api_url', api_url)
        self._save()

    def set_user_update_token(self, user_update_token):
        self.parser.read(self.filename)
        self.parser.set('redirect', 'user_update_token', user_update_token)
        self._save()

    def get_domain(self):
        self.parser.read(self.filename)
        self.logger.info('log file: ' + self.filename)
        return self.parser.get('redirect', 'domain')

    def get_api_url(self):
        self.parser.read(self.filename)
        self.logger.info('log file: ' + self.filename)
        return self.parser.get('redirect', 'api_url')

    def get_user_update_token(self):
        self.parser.read(self.filename)
        return self.parser.get('redirect', 'user_update_token')

    def _save(self):
        self.logger.info('saving config={0}'.format(self.filename))
        with open(self.filename, 'wb') as f:
            self.parser.write(f)
