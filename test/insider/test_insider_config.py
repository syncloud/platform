import unittest

from syncloud.insider.config import InsiderConfig
from test.insider.helpers import temp_file, insider_config


def test_domain():
    filename = temp_file(insider_config)
    config = InsiderConfig(filename)

    config.update('syncloud.it', 'http://api.syncloud.it')
    assert 'syncloud.it' == config.get_redirect_main_domain()
    assert 'http://api.syncloud.it' == config.get_redirect_api_url()

    config.update('syncloud.info', 'http://api.syncloud.info:81')
    assert 'syncloud.info' == config.get_redirect_main_domain()
    assert 'http://api.syncloud.info:81' == config.get_redirect_api_url()