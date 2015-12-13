from os.path import dirname
import convertible
import pytest
import requests
from subprocess import check_output
import time

from integration.util.loop import loop_device_cleanup
from integration.util.ssh import set_docker_ssh_port

from integration.util.ssh import run_ssh

SYNCLOUD_INFO = 'syncloud.info'

DIR = dirname(__file__)
DEVICE_USER = "user1"
DEVICE_PASSWORD = "syncloud"

def test_install(auth):
    __local_install(auth)


def test_non_activated_device_main_page_redirect_to_activation():
    response = requests.get('http://localhost', allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://localhost:81'


def test_non_activated_device_login_redirect_to_activation():
    response = requests.post('http://localhost/server/rest/login', allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://localhost:81'


def test_internal_web_open():

    response = requests.get('http://localhost:81')
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_activate_device(auth):

    email, password, domain, version, arch, release = auth
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': domain, 'name': DEVICE_USER, 'password': DEVICE_PASSWORD,
                                   'api-url': 'http://api.syncloud.info:81', 'domain': SYNCLOUD_INFO})
    assert response.status_code == 200, response.text


def test_reactivate(auth):
    email, password, domain, version, arch, release = auth
    response = requests.post('http://localhost:81/server/rest/activate',
                             data={'redirect-email': email, 'redirect-password': password,
                                   'redirect-domain': domain, 'name': DEVICE_USER, 'password': DEVICE_PASSWORD,
                                   'api-url': 'http://api.syncloud.info:81', 'domain': SYNCLOUD_INFO})
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


def test_default_external_mode_on_activate(auth):

    email, password, domain, version, arch, release = auth

    run_ssh('cp /integration/event/on_domain_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)

    response = session.get('http://localhost/server/rest/settings/external_access')
    assert '"mode": "http"' in response.text
    assert response.status_code == 200

    response = session.get('http://localhost/server/rest/settings/external_access_disable')
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = session.get('http://localhost/server/rest/settings/external_access')
    assert '"mode": null' in response.text
    assert response.status_code == 200
    assert run_ssh('cat /tmp/on_domain_change.log', password=DEVICE_PASSWORD) == '{0}.{1}'.format(domain, SYNCLOUD_INFO)

    response = session.get('http://localhost/server/rest/settings/external_access_enable?mode=http')
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = session.get('http://localhost/server/rest/settings/external_access')
    assert '"mode": "http"' in response.text
    assert response.status_code == 200

    response = session.get('http://localhost/server/rest/settings/external_access_enable?mode=https')
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = session.get('http://localhost/server/rest/settings/external_access')
    assert '"mode": "https"' in response.text
    assert response.status_code == 200


def test_public_web_files():

    response = session.get('http://localhost/server/rest/files')
    assert response.status_code == 200
    response = requests.get('http://localhost/server/rest/files', allow_redirects=False)
    assert response.status_code == 301


def test_do_not_cache_static_files_as_we_get_stale_ui_on_upgrades():

    response = session.get('http://localhost/server/html/settings.html')
    cache_control = response.headers['Cache-Control']
    assert 'no-cache' in cache_control
    assert 'max-age=0' in cache_control


@pytest.yield_fixture(scope='function')
def loop_device():

    loop_device_cleanup(password=DEVICE_PASSWORD)

    print('adding loop device')
    run_ssh('dd if=/dev/zero bs=1M count=10 of=/tmp/disk', password=DEVICE_PASSWORD)
    run_ssh('losetup /dev/loop0 /tmp/disk', password=DEVICE_PASSWORD)
    run_ssh('file -s /dev/loop0', password=DEVICE_PASSWORD)

    yield '/dev/loop0'

    loop_device_cleanup(password=DEVICE_PASSWORD)


def test_public_settings_disk_add_remove_ext4(loop_device):
    __test_fs(loop_device, 'ext4')


def test_public_settings_disk_add_remove_ntfs(loop_device):
    __test_fs(loop_device, 'ntfs')


def __test_fs(loop_device, fs):

    run_ssh('cp /integration/event/on_disk_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)

    run_ssh('mkfs.{0} {1}'.format(fs, loop_device), password=DEVICE_PASSWORD)

    run_ssh('mkdir /tmp/test', password=DEVICE_PASSWORD)

    run_ssh('mount {0} /tmp/test'.format(loop_device), password=DEVICE_PASSWORD)
    for mount in run_ssh('mount', debug=False, password=DEVICE_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh('umount {0}'.format(loop_device), password=DEVICE_PASSWORD)

    run_ssh('udisksctl mount -b {0}'.format(loop_device))
    for mount in run_ssh('mount', debug=False, password=DEVICE_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh('udisksctl unmount -b {0}'.format(loop_device), password=DEVICE_PASSWORD)

    response = session.get('http://localhost/server/rest/settings/disks')
    print response.text
    assert loop_device in response.text
    assert response.status_code == 200

    response = session.get('http://localhost/server/rest/settings/disk_activate',
                           params={'device': loop_device})
    assert response.status_code == 200
    assert run_ssh('cat /tmp/on_disk_change.log', password=DEVICE_PASSWORD) == '/data/platform'

    response = session.get('http://localhost/server/rest/settings/disk_deactivate',
                           params={'device': loop_device})
    assert response.status_code == 200
    assert run_ssh('cat /tmp/on_disk_change.log', password=DEVICE_PASSWORD) == '/data/platform'


def test_internal_web_id():

    response = requests.get('http://localhost:81/server/rest/id')
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_remove():
    run_ssh('/opt/app/sam/bin/sam --debug remove platform', password=DEVICE_PASSWORD)
    time.sleep(3)


def test_reinstall(auth):
    __local_install(auth)


def test_public_web_login_after_reinstall():
    __public_web_login(reset_session=True)


def test_public_web_platform_upgrade():

    response = session.get('http://localhost/server/rest/settings/system_upgrade')
    assert response.status_code == 200
    sam_running = True
    while sam_running:
        try:
            response = session.get('http://localhost/server/rest/settings/sam_status')
            if response.status_code == 200:
                json = convertible.from_json(response.text)
                sam_running = json.is_running
        except Exception, e:
            pass
        time.sleep(1)


def test_reinstall_local_after_upgrade(auth):
    __local_install(auth)


def test_nginx_performance():
    print(check_output('ab -c 1 -n 1000 http://127.0.0.1/ping', shell=True))


def test_nginx_plus_flask_performance():
    print(check_output('ab -c 1 -n 1000 http://127.0.0.1:81/server/rest/id', shell=True))


def __public_web_login(reset_session=False):
    global session
    if reset_session:
        session = requests.session()
    session.post('http://localhost/server/rest/login', data={'name': DEVICE_USER, 'password': DEVICE_PASSWORD})
    assert session.get('http://localhost/server/rest/user', allow_redirects=False).status_code == 200


def __local_install(auth):
    email, password, domain, version, arch, release = auth
    run_ssh('/opt/app/sam/bin/sam --debug install /platform-{0}-{1}.tar.gz'.format(version, arch), password=DEVICE_PASSWORD)
    run_ssh('/opt/app/sam/bin/sam update --release {0}'.format(release), password=DEVICE_PASSWORD)
    set_docker_ssh_port()
    time.sleep(3)


