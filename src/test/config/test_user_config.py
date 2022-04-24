from os.path import dirname, join, isfile

from syncloud_platform.config.user_config import PlatformUserConfig
from test.insider.helpers import temp_file
import os
from os import path


def test_domain():
    config = PlatformUserConfig(temp_file())
    config.init_user_config()

    config.update_redirect('syncloud.it')
    assert 'syncloud.it' == config.get_redirect_domain()
    assert 'https://api.syncloud.it' == config.get_redirect_api_url()

    config.update_redirect('syncloud.info')
    assert 'syncloud.info' == config.get_redirect_domain()
    assert 'https://api.syncloud.info' == config.get_redirect_api_url()


def test_none():
    config_db = join(dirname(__file__), 'db')
    if isfile(config_db):
        os.remove(config_db)
    config = PlatformUserConfig(config_db)
    config.set_web_secret_key(None)
