import os
from os.path import dirname, join
import convertible
import pytest
import requests
from subprocess import check_output
import time

import shutil
import socket
import pytest

from requests.adapters import HTTPAdapter

from integration.util.loop import loop_device_cleanup
from integration.util.ssh import set_docker_ssh_port, run_scp, SSH, ssh_command
from integration.util.ssh import run_ssh
from integration.util.helper import local_install, wait_for_platform_web

SYNCLOUD_INFO = 'syncloud.info'

DIR = dirname(__file__)
DEVICE_USER = "user"
DEVICE_PASSWORD = "password"
DEFAULT_DEVICE_PASSWORD = "syncloud"
LOGS_SSH_PASSWORD = DEFAULT_DEVICE_PASSWORD
LOG_DIR = join(DIR, 'log')


@pytest.fixture(scope="session")
def module_setup(request):
    request.addfinalizer(module_teardown)


def module_teardown():
    os.mkdir(LOG_DIR)
    run_scp('root@localhost:/opt/data/platform/log/* {0}'.format(LOG_DIR), password=LOGS_SSH_PASSWORD)

    print('systemd logs')
    run_ssh('journalctl | grep platform', password=LOGS_SSH_PASSWORD)

    print('-------------------------------------------------------')
    print('syncloud docker image is running')
    print('connect using: {0}'.format(ssh_command(DEVICE_PASSWORD, SSH)))
    print('-------------------------------------------------------')


@pytest.fixture(scope="function")
def public_web_session():
    wait_for_platform_web()
    session = requests.session()
    session.post('http://localhost/rest/login', data={'name': DEVICE_USER, 'password': DEVICE_PASSWORD})
    assert session.get('http://localhost/rest/user', allow_redirects=False).status_code == 200
    return session


@pytest.fixture(scope="module")
def user_domain(auth):
    email, password, domain, version, arch, release = auth
    return '{0}.{1}'.format(domain, SYNCLOUD_INFO)

def test_start(module_setup):
    shutil.rmtree(LOG_DIR, ignore_errors=True)


def test_install(auth):
    email, password, domain, version, arch, release = auth
    local_install(DIR, DEFAULT_DEVICE_PASSWORD, version, arch, release)


def test_activate_device(auth):

    email, password, domain, version, arch, release = auth
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = 'password'
    response = requests.post('http://localhost:81/rest/activate',
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': 'user', 'device_password': 'password'})
    assert response.status_code == 200, response.text


def test_external_mode(auth, public_web_session, user_domain):

    email, password, domain, version, arch, release = auth

    response = public_web_session.get('http://localhost/rest/settings/set_external_access',
                                      params={'external_access': 'true'})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/external_access')
    assert '"external_access": true' in response.text
    assert response.status_code == 200

    _wait_for_ip(user_domain)


def _wait_for_ip(user_domain):

    retries = 10
    retry = 0
    while retry < retries:
        ip = socket.gethostbyname(user_domain)
        if not ip.startswith('192'):
            return
        retry += 1
        time.sleep(1)

def test_certbot_cli():
    run_ssh('/opt/app/platform/bin/certbot --help', password=DEVICE_PASSWORD)


def test_external_https_mode_with_certbot(public_web_session):

    response = public_web_session.get('http://localhost/rest/settings/regenerate_certificate')
    assert '"success": true' in response.text
    assert response.status_code == 200