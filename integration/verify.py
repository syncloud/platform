import os
from os.path import dirname, join
import convertible
import requests
from subprocess import check_output
import time

import shutil
import socket
import pytest

from requests.adapters import HTTPAdapter

from integration.util.loop import loop_device_cleanup
from integration.util.ssh import run_scp, SSH, ssh_command
from integration.util.ssh import run_ssh
from integration.util.helper import local_install, wait_for_platform_web, wait_for_sam

SYNCLOUD_INFO = 'syncloud.info'

DIR = dirname(__file__)
DEVICE_USER = "user"
DEVICE_PASSWORD = "password"
DEFAULT_DEVICE_PASSWORD = "syncloud"
LOGS_SSH_PASSWORD = DEFAULT_DEVICE_PASSWORD
LOG_DIR = join(DIR, 'log')

SAM_DATA_DIR='/opt/data/platform'
SNAPD_DATA_DIR='/var/snap/platform/common'
DATA_DIR=''

SAM_APP_DIR='/opt/app/platform'
SNAPD_APP_DIR='/snap/platform/current'
APP_DIR=''

@pytest.fixture(scope="session")
def data_dir(installer):
    if installer == 'sam':
        return SAM_DATA_DIR
    else:
        return SNAPD_DATA_DIR


@pytest.fixture(scope="session")
def app_dir(installer):
    if installer == 'sam':
        return SAM_APP_DIR
    else:
        return SNAPD_APP_DIR


@pytest.fixture(scope="session")
def service_prefix(installer):
    if installer == 'sam':
        return ''
    else:
        return 'snap.'


@pytest.fixture(scope="session")
def conf_dir(installer):
    if installer == 'sam':
        return SAM_APP_DIR
    else:
        return SNAPD_DATA_DIR


@pytest.fixture(scope="session")
def module_setup(request, data_dir):
    global DATA_DIR
    DATA_DIR=data_dir
    request.addfinalizer(module_teardown)


def module_teardown():
    os.mkdir(LOG_DIR)
    run_scp('root@localhost:{0}/log/* {1}'.format(DATA_DIR, LOG_DIR), password=LOGS_SSH_PASSWORD)

    print('systemd logs')
    run_ssh('journalctl | tail -200', password=LOGS_SSH_PASSWORD)

    print('-------------------------------------------------------')
    print('syncloud docker image is running')
    print('connect using: {0}'.format(ssh_command(DEVICE_PASSWORD, SSH)))
    print('-------------------------------------------------------')


def test_start(module_setup):
    shutil.rmtree(LOG_DIR, ignore_errors=True)


def test_install(auth, installer):
    email, password, domain, app_archive_path = auth
    local_install(DEFAULT_DEVICE_PASSWORD, app_archive_path, installer)


def test_non_activated_device_main_page_redirect_to_activation():
    response = requests.get('http://localhost', allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://localhost:81'


def test_non_activated_device_login_redirect_to_activation():
    response = requests.post('http://localhost/rest/login', allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://localhost:81'


def test_internal_web_open():

    response = requests.get('http://localhost:81')
    assert response.status_code == 200


def test_activate_device(auth):

    email, password, domain, app_archive_path = auth
    global LOGS_SSH_PASSWORD
    response = requests.post('http://localhost:81/rest/activate',
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': 'user1', 'device_password': 'password1'})
    assert response.status_code == 200, response.text
    LOGS_SSH_PASSWORD = 'password1'


def test_reactivate(auth):
    email, password, domain, app_archive_path = auth
    response = requests.post('http://localhost:81/rest/activate',
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': DEVICE_USER, 'device_password': DEVICE_PASSWORD})
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = DEVICE_PASSWORD


def test_public_web_unauthorized_browser_redirect():
    response = requests.get('http://localhost/rest/user', allow_redirects=False)
    assert response.status_code == 302


def test_public_web_unauthorized_ajax_not_redirect():
    response = requests.get('http://localhost/rest/user',
                            allow_redirects=False, headers={'X-Requested-With': 'XMLHttpRequest'})
    assert response.status_code == 401


def test_running_platform_web():
    print(check_output('nc -zv -w 1 localhost 80', shell=True))


def test_platform_rest():
    session = requests.session()
    session.mount('http://localhost', HTTPAdapter(max_retries=5))
    response = session.get('http://localhost', timeout=60)
    assert response.status_code == 200


# def test_external_mode(auth, public_web_session, user_domain):
#
#     email, password, domain, app_archive_path = auth
#
#     run_ssh('cp /integration/event/on_domain_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)
#
#     response = public_web_session.get('http://localhost/rest/settings/external_access')
#     assert '"external_access": false' in response.text
#     assert response.status_code == 200
#
#     response = public_web_session.get('http://localhost/rest/settings/set_external_access',
#                                       params={'external_access': 'true'})
#     assert '"success": true' in response.text
#     assert response.status_code == 200
#
#     response = public_web_session.get('http://localhost/rest/settings/external_access')
#     assert '"external_access": true' in response.text
#     assert response.status_code == 200
#
#     _wait_for_ip(user_domain)
#
#     assert run_ssh('cat /tmp/on_domain_change.log', password=DEVICE_PASSWORD) == '{0}.{1}'.format(domain, SYNCLOUD_INFO)


def _wait_for_ip(user_domain):

    retries = 10
    retry = 0
    while retry < retries:
        ip = socket.gethostbyname(user_domain)
        if not ip.startswith('192'):
            return
        retry += 1
        time.sleep(1)

def test_certbot_cli(app_dir):
    run_ssh('{0}/bin/certbot --help'.format(app_dir), password=DEVICE_PASSWORD)


def test_external_https_mode_with_certbot(public_web_session):

    response = public_web_session.get('http://localhost/rest/settings/set_protocol',
                                      params={'protocol': 'https'})
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_show_https_certificate():
    run_ssh("echo | "
            "openssl s_client -showcerts -servername localhost -connect localhost:443 2>/dev/null | "
            "openssl x509 -inform pem -noout -text", password=DEVICE_PASSWORD)


def test_protocol(auth, public_web_session, conf_dir, service_prefix):

    email, password, domain, app_archive_path = auth

    run_ssh("sed -i 's#hooks_root.*#hooks_root: /integration#g' {0}/config/platform.cfg".format(conf_dir), password=DEVICE_PASSWORD)

    run_ssh('systemctl restart {0}platform.uwsgi-public'.format(service_prefix), password=DEVICE_PASSWORD)

    wait_for_platform_web()

    response = public_web_session.get('http://localhost/rest/settings/protocol')
    assert '"protocol": "https"' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/set_protocol',
                                      params={'protocol': 'https'})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/protocol')
    assert '"protocol": "https"' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/set_protocol',
                                      params={'protocol': 'http'})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/protocol')
    assert '"protocol": "http"' in response.text
    assert response.status_code == 200

    assert run_ssh('cat /tmp/on_domain_change.log', password=DEVICE_PASSWORD) == '{0}.{1}'.format(domain, SYNCLOUD_INFO)


def test_cron_job(auth, public_web_session, app_dir):
    assert '"success": true' in run_ssh('{0}/bin/insider sync_all'.format(app_dir), password=DEVICE_PASSWORD)


def test_installed_apps(public_web_session):
    response = public_web_session.get('http://localhost/rest/installed_apps')
    assert response.status_code == 200


def test_do_not_cache_static_files_as_we_get_stale_ui_on_upgrades(public_web_session):

    response = public_web_session.get('http://localhost/settings.html')
    cache_control = response.headers['Cache-Control']
    assert 'no-cache' in cache_control
    assert 'max-age=0' in cache_control


def test_installer_upgrade(public_web_session, installer):
    __upgrade(public_web_session, 'sam')


@pytest.yield_fixture(scope='function')
def loop_device():
    dev_file = '/tmp/disk'
    loop_device_cleanup(dev_file, password=DEVICE_PASSWORD)

    print('adding loop device')
    run_ssh('dd if=/dev/zero bs=1M count=10 of={0}'.format(dev_file), password=DEVICE_PASSWORD)
    loop = run_ssh('losetup -f --show {0}'.format(dev_file), password=DEVICE_PASSWORD)
    run_ssh('file -s {0}'.format(loop), password=DEVICE_PASSWORD)

    yield loop

    loop_device_cleanup(dev_file, password=DEVICE_PASSWORD)


def disk_writable():
    run_ssh('ls -la /data/', password=DEVICE_PASSWORD)
    run_ssh("su - platform -s /bin/bash -c 'touch /data/platform/test.file'", password=DEVICE_PASSWORD)


@pytest.mark.parametrize("fs_type", ['ext2', 'ext3', 'ext4'])
def test_public_settings_disk_add_remove(loop_device, public_web_session, fs_type):
    disk_create(loop_device, fs_type)
    assert disk_activate(loop_device,  public_web_session) == '/opt/disk/external/platform'
    disk_writable()
    assert disk_deactivate(loop_device, public_web_session) == '/opt/disk/internal/platform'


def test_disk_physical_remove(loop_device, public_web_session):
    disk_create(loop_device, 'ext4')
    assert disk_activate(loop_device,  public_web_session) == '/opt/disk/external/platform'
    loop_device_cleanup('/opt/disk/external', password=DEVICE_PASSWORD)
    run_ssh('udevadm trigger --action=remove -y {0}'.format(loop_device.split('/')[2]), password=DEVICE_PASSWORD)
    run_ssh('udevadm settle', password=DEVICE_PASSWORD)
    assert current_disk_link() == '/opt/disk/internal/platform'


def disk_create(loop_device, fs):
    run_ssh('mkfs.{0} {1}'.format(fs, loop_device), password=DEVICE_PASSWORD)

    run_ssh('rm -rf /tmp/test', password=DEVICE_PASSWORD)
    run_ssh('mkdir /tmp/test', password=DEVICE_PASSWORD)

    run_ssh('mount {0} /tmp/test'.format(loop_device), password=DEVICE_PASSWORD)
    for mount in run_ssh('mount', debug=True, password=DEVICE_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh('umount {0}'.format(loop_device), password=DEVICE_PASSWORD)


def disk_activate(loop_device, public_web_session):

    run_ssh('cp /integration/event/on_disk_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)

    response = public_web_session.get('http://localhost/rest/settings/disks')
    print response.text
    assert loop_device in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/disk_activate',
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link()


def disk_deactivate(loop_device, public_web_session):
    response = public_web_session.get('http://localhost/rest/settings/disk_deactivate',
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link()


def current_disk_link():
    return run_ssh('cat /tmp/on_disk_change.log', password=DEVICE_PASSWORD)


def test_internal_web_id():

    response = requests.get('http://localhost:81/rest/id')
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_if_cron_is_enabled_after_install():
    cron_is_enabled_after_install()


def cron_is_enabled_after_install():
    crontab = run_ssh("crontab -l", password=DEVICE_PASSWORD)
    assert len(crontab.splitlines()) == 1
    assert 'cron.py' in crontab, crontab
    assert not crontab.startswith('#'), crontab


def test_remove():
    run_ssh('/opt/app/sam/bin/sam --debug remove platform', password=DEVICE_PASSWORD)
    time.sleep(3)


def test_reinstall(auth, installer):
    email, password, domain, app_archive_path = auth
    local_install(DEVICE_PASSWORD, app_archive_path, installer)


def test_public_web_platform_upgrade(public_web_session):
    __upgrade(public_web_session, 'system')


def __upgrade(public_web_session, upgrade_type):

    public_web_session.get('http://localhost/rest/settings/{0}_upgrade'.format(upgrade_type))
    wait_for_sam(public_web_session)


def test_reinstall_local_after_upgrade(auth, installer):
    email, password, domain, app_archive_path = auth
    local_install(DEVICE_PASSWORD, app_archive_path, installer)


def test_if_cron_is_enabled_after_upgrade():
    cron_is_enabled_after_install()


def test_nginx_performance():
    print(check_output('ab -c 1 -n 1000 http://127.0.0.1/ping', shell=True))


def test_nginx_plus_flask_performance():
    print(check_output('ab -c 1 -n 1000 http://127.0.0.1:81/rest/id', shell=True))
