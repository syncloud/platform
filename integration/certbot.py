import os
import shutil
import socket
import time
from os.path import dirname, join
from subprocess import check_output
import pytest
import requests
from requests.adapters import HTTPAdapter
from integration.util.helper import local_install, wait_for_platform_web
from integration.util.ssh import run_scp, SSH, ssh_command
from integration.util.ssh import run_ssh

SYNCLOUD_INFO = 'syncloud.info'

DIR = dirname(__file__)
DEVICE_USER = "user"
DEVICE_PASSWORD = "password"
DEFAULT_DEVICE_PASSWORD = "syncloud"
LOGS_SSH_PASSWORD = DEFAULT_DEVICE_PASSWORD
LOG_DIR = join(DIR, 'log')
CERTIFICATE_DIR = join(DIR, 'certificate')


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
    email, password, domain, app_archive_path = auth
    local_install(DEFAULT_DEVICE_PASSWORD, app_archive_path)


def test_activate_device(auth):

    email, password, domain, app_archive_path = auth
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = 'password'
    response = requests.post('http://localhost:81/rest/activate',
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': 'user', 'device_password': 'password'})
    assert response.status_code == 200, response.text


def test_running_platform_web():
    print(check_output('nc -zv -w 1 localhost 80', shell=True))


def test_platform_rest():
    session = requests.session()
    session.mount('http://localhost', HTTPAdapter(max_retries=5))
    response = session.get('http://localhost', timeout=60)
    assert response.status_code == 200


def test_external_mode(public_web_session, user_domain):

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

    os.mkdir(CERTIFICATE_DIR)

    run_scp('root@localhost:/opt/data/platform/certbot/live/*/fullchain.pem {0}/syncloud.crt'
            .format(CERTIFICATE_DIR),
            password=LOGS_SSH_PASSWORD)

    run_scp('root@localhost:/opt/data/platform/certbot/keys/0000_key-certbot.pem {0}/syncloud.key'
            .format(CERTIFICATE_DIR),
            password=LOGS_SSH_PASSWORD)
