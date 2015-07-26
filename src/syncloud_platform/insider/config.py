from ConfigParser import ConfigParser
from os.path import dirname, join
import os
from syncloud_app import logger
import convertible


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


class RedirectConfig:
    def __init__(self, filename):
        self.parser = ConfigParser()
        self.parser.read(filename)
        self.filename = filename
        self.logger = logger.get_logger('insider.RedirectConfig')

    def update(self, domain, api_url):
        self.logger.info('settig domain={0}, api_url={1}'.format(domain, api_url))
        if not self.parser.has_section('redirect'):
            self.parser.add_section('redirect')

        self.parser.set('redirect', 'domain', domain)
        self.parser.set('redirect', 'api_url', api_url)
        self._save()

    def get_domain(self):
        return self.parser.get('redirect', 'domain')

    def get_api_url(self):
        return self.parser.get('redirect', 'api_url')

    def _save(self):
        self.logger.info('saving config={0}'.format(self.filename))
        with open(self.filename, 'wb') as file:
            self.parser.write(file)

INSIDER_CONFIG_NAME = 'insider.cfg'


class InsiderConfig:
    # TODO: Split redirect and insider config files
    def __init__(self, config_dir, redirect_config):
        self.parser = ConfigParser()
        self.filename = join(config_dir, INSIDER_CONFIG_NAME)
        self.parser.read(self.filename)
        self.redirect_config = redirect_config
        self.logger = logger.get_logger('insider.InsiderConfig')

    def _save(self):
        self.logger.info('saving config={0}'.format(self.filename))
        with open(self.filename, 'wb') as f:
            self.parser.write(f)

    def update(self, domain, api_url):
        self.redirect_config.update(domain, api_url)

    def get_redirect_api_url(self):
        return self.redirect_config.get_api_url()

    def get_redirect_main_domain(self):
        return self.redirect_config.get_domain()

    def get_cron_period_mins(self):
        return self.parser.getint('insider', 'cron_period_mins')

    def get_upnp_enabled(self):
        return self.parser.getboolean('insider', 'upnp_enabled')

    def set_upnp_enabled(self, enabled):
        self.parser.set('insider', 'upnp_enabled', str(enabled))
        self._save()
