from os.path import dirname
import convertible

import requests
from subprocess import check_output
import time

DIR = dirname(__file__)


def test_install(auth):
    __local_install(auth)


def test_non_activated_device_redirect_to_activation():
    response = requests.post('http://localhost/server/rest/login', allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://localhost:81'


def test_internal_web_open():

    response = requests.get('http://localhost:81')
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_activate_device(auth):

    email, password, domain, version, arch = auth
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': domain, 'name': 'user1', 'password': 'password1',
                                   'api-url': 'http://api.syncloud.info:81', 'domain': 'syncloud.info'})
    assert response.status_code == 200, response.text


def test_reactivate(auth):
    email, password, domain, version, arch = auth
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': domain, 'name': 'user', 'password': 'password',
                                   'api-url': 'http://api.syncloud.info:81', 'domain': 'syncloud.info'})
    assert response.status_code == 200


def test_public_web_unauthorized_browser_redirect():
    response = requests.get('http://localhost/server/rest/user', allow_redirects=False)
    assert response.status_code == 302


def test_public_web_unauthorized_ajax_not_redirect():
    response = requests.get('http://localhost/server/rest/user',
                            allow_redirects=False, headers={'X-Requested-With': 'XMLHttpRequest'})
    assert response.status_code == 401

session = requests.session()


def test_public_web_login():
    __public_web_login()


def test_public_web_files():

    response = session.get('http://localhost/server/rest/files')
    assert response.status_code == 200
    response = requests.get('http://localhost/server/rest/files', allow_redirects=False)
    assert response.status_code == 301


def test_public_settings_disk_add_remove():

    # device = check_output('{0}/virtual_disk.sh add'.format(DIR), shell=True).strip()

    response = session.get('http://localhost/server/rest/settings/disks')
    assert response.status_code == 200
    # print response.text
    # print response.url
    json = convertible.from_json(response.text)
    for disk in json.disks:
        for partition in disk.partitions:
            response = session.get('http://localhost/server/rest/settings/disk_activate', params={
                'device': partition.device, 'fix_permissions': False})
            assert response.status_code == 200
            response = session.get(
                'http://localhost/server/rest/settings/disk_deactivate', params={'device': partition.device})
            assert response.status_code == 200

    # check_output('{0}/virtual_disk.sh remove'.format(DIR), shell=True)


def test_internal_web_id():

    response = requests.get('http://localhost:81/server/rest/id')
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_remove():
    ssh = 'sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost'
    print(check_output('{0} /opt/app/sam/bin/sam --debug remove platform'.format(ssh), shell=True))
    time.sleep(3)


def test_reinstall(auth):
    __local_install(auth)


def test_public_web_login_after_reinstall():
    __public_web_login(reset_session=True)


def __public_web_login(reset_session=False):
    global session
    if reset_session:
        session = requests.session()
    session.post('http://localhost/server/rest/login', data={'name': 'user', 'password': 'password'})
    assert session.get('http://localhost/server/rest/user', allow_redirects=False).status_code == 200


def __local_install(auth):
    email, password, domain, version, arch = auth
    ssh = 'sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p 2222 root@localhost'
    print(check_output('{0} /opt/app/sam/bin/sam --debug install /platform-{1}-{2}.tar.gz'.format(ssh, version, arch),
                       shell=True))
    time.sleep(3)
