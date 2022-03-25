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

def test_migrate():
    old_config_file = temp_file()
    with open(old_config_file, 'w') as f:
       f.write("""
[platform]
redirect_enabled = True
user_domain = test
domain_update_token = token1
external_access = False
manual_access_port = 443
activated = True

[redirect]
domain = syncloud.it
api_url = http://api.syncloud.it
user_email = user@example.com
user_update_token = token2
       """)
    config_db = join(dirname(__file__), 'db')
    if isfile(config_db):
        os.remove(config_db)
    config = PlatformUserConfig(config_db, old_config_file)
    assert config.get_redirect_domain() == 'syncloud.it'
    assert config.get_upnp()
    assert config.is_redirect_enabled()
    assert not config.get_external_access()

    assert not path.isfile(old_config_file)


def test_none():
    config_db = join(dirname(__file__), 'db')
    if isfile(config_db):
        os.remove(config_db)
    config = PlatformUserConfig(config_db)
    config.set_web_secret_key(None)
