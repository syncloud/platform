import logging
from os.path import dirname
import requests
import pytest
from subprocess import check_output
from syncloud.sam.platform_installer import PlatformInstaller
from syncloud.server.serverfacade import get_server
from syncloud.insider.facade import get_insider
from syncloud.sam.pip import Pip
from syncloud.app import logger

DIR = dirname(__file__)

@pytest.fixture(scope="session", autouse=True)
def activate_device(auth):

    logger.init(logging.DEBUG, True)

    print("installing local binary build")
    PlatformInstaller().install('platform.tar.gz')

    Pip(None).log_version('syncloud-platform')

    # persist upnp mock setting
    get_insider().insider_config.set_upnpc_mock(True)

    server = get_server(insider=get_insider(use_upnpc_mock=True))
    email, password = auth
    server.activate('test', 'syncloud.info', 'http://api.syncloud.info:81', email, password, 'teamcity', 'user', 'password', False)

    # request.addfinalizer(finalizer_function)

def test_public_web():
    session = requests.session()
    response = session.get('http://localhost/server/rest/user', allow_redirects=False)
    print(response.text)
    assert response.status_code == 302
    response = session.post('http://localhost/server/rest/login', data={'name': 'user', 'password': 'password'})
    print(response.text)
    assert session.get('http://localhost/server/rest/user', allow_redirects=False).status_code == 200

def test_internal_web():

    response = requests.get('http://localhost:81/id')
    assert 'mac_address' in response.text
    assert response.status_code == 200
