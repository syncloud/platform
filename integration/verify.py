import logging
from os.path import dirname

import requests
import pytest

from syncloud.app import logger

DIR = dirname(__file__)

@pytest.fixture(scope="session", autouse=True)
def activate_device(auth):

    logger.init(logging.DEBUG, True)

    email, password = auth

    release = open('{0}/RELEASE'.format(DIR), 'r').read().strip()

    # activate
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': 'teamcity', 'name': 'user1', 'password': 'password1',
                                   'api-url': 'http://api.syncloud.info:81', 'domain': 'syncloud.info',
                                   'release': release})
    assert response.status_code == 200

    # re-activate
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': 'teamcity', 'name': 'user', 'password': 'password',
                                   'api-url': 'http://api.syncloud.info:81', 'domain': 'syncloud.info',
                                   'release': release})
    assert response.status_code == 200


def test_public_web():
    session = requests.session()
    response = session.get('http://localhost/server/rest/user', allow_redirects=False)
    print(response.text)
    assert response.status_code == 302
    response = session.post('http://localhost/server/rest/login', data={'name': 'user', 'password': 'password'})
    print(response.text)
    assert session.get('http://localhost/server/rest/user', allow_redirects=False).status_code == 200

def test_internal_web():

    response = requests.get('http://localhost:81/server/rest/id')
    assert 'mac_address' in response.text
    assert response.status_code == 200
