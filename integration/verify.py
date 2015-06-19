import logging
from os.path import dirname

import requests
import pytest

from syncloud.app import logger

DIR = dirname(__file__)

def test_non_activated_device_redirect_to_activation():
    response = requests.post('http://localhost/server/rest/login', allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://localhost:81'

def test_activate_device(auth):

    email, password, domain = auth
    release = open('{0}/RELEASE'.format(DIR), 'r').read().strip()
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': domain, 'name': 'user1', 'password': 'password1',
                                   'api-url': 'http://api.syncloud.info:81', 'domain': 'syncloud.info',
                                   'release': release})
    assert response.status_code == 200

def test_reactivate(auth):
    email, password, domain = auth
    release = open('{0}/RELEASE'.format(DIR), 'r').read().strip()
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': domain, 'name': 'user', 'password': 'password',
                                   'api-url': 'http://api.syncloud.info:81', 'domain': 'syncloud.info',
                                   'release': release})
    assert response.status_code == 200

def test_public_web_unauthorized_browser_redirect():
    response = requests.get('http://localhost/server/rest/user', allow_redirects=False)
    assert response.status_code == 302

def test_public_web_unauthorized_ajax_not_redirect():
    response = requests.get('http://localhost/server/rest/user', allow_redirects=False, headers={'X-Requested-With': 'XMLHttpRequest'})
    assert response.status_code == 401

def test_public_web_login():
    session = requests.session()
    session.post('http://localhost/server/rest/login', data={'name': 'user', 'password': 'password'})
    assert session.get('http://localhost/server/rest/user', allow_redirects=False).status_code == 200

def test_internal_web_open():

    response = requests.get('http://localhost:81/server/rest/id')
    assert 'mac_address' in response.text
    assert response.status_code == 200
