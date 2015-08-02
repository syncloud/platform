from os.path import dirname

from syncloud_platform.insider.config import RedirectConfig, REDIRECT_CONFIG_NAME
from test.insider.helpers import temp_file


def test_domain():
    config = RedirectConfig(dirname(temp_file(filename=REDIRECT_CONFIG_NAME)))

    config.update('syncloud.it', 'http://api.syncloud.it')
    assert 'syncloud.it' == config.get_domain()
    assert 'http://api.syncloud.it' == config.get_api_url()

    config.update('syncloud.info', 'http://api.syncloud.info:81')
    assert 'syncloud.info' == config.get_domain()
    assert 'http://api.syncloud.info:81' == config.get_api_url()
