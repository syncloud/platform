import unittest

from syncloud.insider.config import InsiderConfig
from test.insider.helpers import temp_file, insider_config


class TestInsiderConfig(unittest.TestCase):

    def test_domain(self):
        filename = temp_file(insider_config)
        config = InsiderConfig(filename)

        config.update('syncloud.it', 'http://api.syncloud.it')
        self.assertEquals('syncloud.it', config.get_redirect_main_domain())
        self.assertEquals('http://api.syncloud.it', config.get_redirect_api_url())

        config.update('syncloud.info', 'http://api.syncloud.info:81')
        self.assertEquals('syncloud.info', config.get_redirect_main_domain())
        self.assertEquals('http://api.syncloud.info:81', config.get_redirect_api_url())