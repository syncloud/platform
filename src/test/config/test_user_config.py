from syncloud_platform.config.user_config import PlatformUserConfig
from test.insider.helpers import temp_file
from os import path

def test_domain():
    config = PlatformUserConfig(temp_file())
    config.init_user_config()
    
    config.update_redirect('syncloud.it', 'http://api.syncloud.it')
    assert 'syncloud.it' == config.get_redirect_domain()
    assert 'http://api.syncloud.it' == config.get_redirect_api_url()

    config.update_redirect('syncloud.info', 'http://api.syncloud.info:81')
    assert 'syncloud.info' == config.get_redirect_domain()
    assert 'http://api.syncloud.info:81' == config.get_redirect_api_url()

def test_migrate():
    old_config_file = temp_file()
    with open(old_config_file, 'w') as f:
       f.write("""
[platform]
redirect_enabled = True
user_domain = test
domain_update_token = token1
external_access = False
manual_certificate_port = 80
manual_access_port = 443
activated = True

[redirect]
domain = syncloud.it
api_url = http://api.syncloud.it
user_email = user@example.com
user_update_token = token2
       """)
    config = PlatformUserConfig(temp_file(), old_config_file)

    assert config.get_redirect_domain() == 'syncloud.it'
    assert config.get_upnp() == True
    assert config.is_redirect_enabled() == True
    assert config.get_external_access() == False
    
    assert not path.isfile(old_config_file)

def test_none():
    config = PlatformUserConfig(temp_file())

    config.set_web_secret_key(None)
    
