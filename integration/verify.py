import os
from os.path import dirname, join
import convertible
import pytest
import requests
from subprocess import check_output
import time

import shutil

import pytest

from requests.adapters import HTTPAdapter

from integration.util.loop import loop_device_cleanup
from integration.util.ssh import set_docker_ssh_port, run_scp, SSH, ssh_command
from integration.util.ssh import run_ssh

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


def test_start(module_setup):
    shutil.rmtree(LOG_DIR, ignore_errors=True)


def test_install(auth):
    email, password, domain, version, arch, release = auth
    __local_install(DEFAULT_DEVICE_PASSWORD, version, arch, release)


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
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_activate_device(auth):

    email, password, domain, version, arch, release = auth
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = 'password1'
    response = requests.post('http://localhost:81/rest/activate',
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': 'user1', 'device_password': 'password1'})
    assert response.status_code == 200, response.text


def test_reactivate(auth):
    email, password, domain, version, arch, release = auth
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


def test_external_https_mode_with_certbot(public_web_session):

    run_ssh("sed -i 's/certbot_enabled:.*/certbot_enabled: True/g' /opt/app/platform/config/platform.cfg", password=DEVICE_PASSWORD)

    response = public_web_session.get('http://localhost/rest/settings/set_external_access',
                                      params={'external_access': 'true'})
    assert '"success": true' in response.text
    assert response.status_code == 200
    
    response = public_web_session.get('http://localhost/rest/settings/set_protocol',
                                      params={'protocol': 'https'})
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_external_mode(auth, public_web_session):

    email, password, domain, version, arch, release = auth

    run_ssh('cp /integration/event/on_domain_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)

    response = public_web_session.get('http://localhost/rest/settings/external_access')
    assert '"external_access": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/set_external_access',
                                      params={'external_access': 'false'})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/external_access')
    assert '"external_access": false' in response.text
    assert response.status_code == 200

    assert run_ssh('cat /tmp/on_domain_change.log', password=DEVICE_PASSWORD) == '{0}.{1}'.format(domain, SYNCLOUD_INFO)


def test_protocol(auth, public_web_session):

    email, password, domain, version, arch, release = auth

    run_ssh('cp /integration/event/on_domain_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)

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


def test_cron_job(auth, public_web_session):
    assert '"success": true' in run_ssh('/opt/app/platform/bin/insider sync_all', password=DEVICE_PASSWORD)


def test_installed_apps(public_web_session):
    response = public_web_session.get('http://localhost/rest/installed_apps')
    assert response.status_code == 200


def test_do_not_cache_static_files_as_we_get_stale_ui_on_upgrades(public_web_session):

    response = public_web_session.get('http://localhost/settings.html')
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


def disk_writable():
    run_ssh('ls -la /data/', password=DEVICE_PASSWORD)
    run_ssh("su - platform -s /bin/bash -c 'touch /data/platform/test.file'", password=DEVICE_PASSWORD)


@pytest.mark.parametrize("fs_type", ['ntfs', 'vfat', 'exfat', 'ext2', 'ext3', 'ext4'])
def test_public_settings_disk_add_remove(loop_device, public_web_session, fs_type):
    disk_create(loop_device, fs_type)
    assert disk_activate(loop_device,  public_web_session) == '/opt/disk/external/platform'
    disk_writable()
    assert disk_deactivate(loop_device, public_web_session) == '/opt/disk/internal/platform'


def test_disk_physical_remove(loop_device, public_web_session):
    disk_create(loop_device, 'ext4')
    assert disk_activate(loop_device,  public_web_session) == '/opt/disk/external/platform'
    loop_device_cleanup(password=DEVICE_PASSWORD)
    run_ssh('udevadm trigger --action=remove -y loop0', password=DEVICE_PASSWORD)
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

    run_ssh('udisksctl mount -b {0}'.format(loop_device), password=DEVICE_PASSWORD)
    for mount in run_ssh('mount', debug=True, password=DEVICE_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh('udisksctl unmount -b {0}'.format(loop_device), password=DEVICE_PASSWORD)


def disk_activate(loop_device,  public_web_session):

    run_ssh('cp /integration/event/on_disk_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)

    response = public_web_session.get('http://localhost/rest/settings/disks')
    print response.text
    assert loop_device in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://localhost/rest/settings/disk_activate',
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link()


def disk_deactivate(loop_device,  public_web_session):
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


def test_reinstall(auth):
    email, password, domain, version, arch, release = auth
    __local_install(DEVICE_PASSWORD, version, arch, release)


def test_public_web_platform_upgrade(public_web_session):

    public_web_session.get('http://localhost/rest/settings/system_upgrade')
    sam_running = True
    while sam_running:
        try:
            response = public_web_session.get('http://localhost/rest/settings/sam_status')
            if response.status_code == 200:
                json = convertible.from_json(response.text)
                sam_running = json.is_running
        except Exception, e:
            pass
        time.sleep(1)


def test_reinstall_local_after_upgrade(auth):
    email, password, domain, version, arch, release = auth
    __local_install(DEVICE_PASSWORD, version, arch, release)


def test_if_cron_is_enabled_after_upgrade():
    cron_is_enabled_after_install()


def test_nginx_performance():
    print(check_output('ab -c 1 -n 1000 http://127.0.0.1/ping', shell=True))


def test_nginx_plus_flask_performance():
    print(check_output('ab -c 1 -n 1000 http://127.0.0.1:81/rest/id', shell=True))


def __local_install(password, version, arch, release):
    run_scp('{0}/../platform-{1}-{2}.tar.gz root@localhost:/'.format(DIR, version, arch), password=password)
    run_ssh('/opt/app/sam/bin/sam --debug install /platform-{0}-{1}.tar.gz'.format(version, arch), password=password)
    run_ssh('/opt/app/sam/bin/sam update --release {0}'.format(release), password=password)
    set_docker_ssh_port(password)
    run_ssh('systemctl restart platform-uwsgi-public', password=password)
    time.sleep(3)


def wait_for_platform_web():
    print(check_output('while ! nc -w 1 -z localhost 81; do sleep 1; done', shell=True))
    print(check_output('while ! nc -w 1 -z localhost 80; do sleep 1; done', shell=True))
