import syncloud_platform.importlib

import unittest

from syncloud_platform.insider.config import InsiderConfig, RedirectConfig
from test.insider.helpers import temp_file, insider_config


def test_domain():
    config = RedirectConfig(temp_file())

    config.update('syncloud.it', 'http://api.syncloud.it')
    assert 'syncloud.it' == config.get_domain()
    assert 'http://api.syncloud.it' == config.get_api_url()

    config.update('syncloud.info', 'http://api.syncloud.info:81')
    assert 'syncloud.info' == config.get_domain()
    assert 'http://api.syncloud.info:81' == config.get_api_url()
