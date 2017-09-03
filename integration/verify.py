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
from integration.util.ssh import run_scp, ssh_command
from integration.util.ssh import run_ssh
from integration.util.helper import local_install, wait_for_sam, wait_for_rest, local_remove

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
        os.environ['SNAP_COMMON'] = SNAPD_DATA_DIR
        return SNAPD_DATA_DIR


@pytest.fixture(scope="session")
def ssh_env_vars(installer):
    if installer == 'sam':
        return ''
    if installer == 'snapd':
        return 'SNAP_COMMON={0} '.format(SNAPD_DATA_DIR)


@pytest.fixture(scope="session")
def module_setup(request, data_dir, device_host, conf_dir):
    request.addfinalizer(lambda: module_teardown(data_dir, device_host, conf_dir))


def module_teardown(data_dir, device_host, conf_dir):
    run_scp('root@{0}:{1}/log/* {2}'.format(device_host, data_dir, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)
    run_scp('-r root@{0}:{1}/config {2}'.format(device_host, conf_dir, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)
    run_scp('root@{0}:/var/log/sam.log {1}'.format(device_host, data_dir, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)

    print('systemd logs')
    run_ssh(device_host, 'journalctl | tail -200', password=LOGS_SSH_PASSWORD)


def test_start(module_setup, device_host):
    shutil.rmtree(LOG_DIR, ignore_errors=True)
    run_scp('-r {0} root@{1}:/'.format(DIR, device_host))
    os.mkdir(LOG_DIR)


def test_install(app_archive_path, installer, device_host):
    run_ssh(device_host, 'systemctl', password=LOGS_SSH_PASSWORD)

    local_install(device_host, DEFAULT_DEVICE_PASSWORD, app_archive_path, installer)


def test_non_activated_device_main_page_redirect_to_activation(device_host):
    response = requests.get('http://{0}'.format(device_host), allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://{0}:81'.format(device_host)


def test_non_activated_device_login_redirect_to_activation(device_host):
    response = requests.post('http://{0}/rest/login'.format(device_host), allow_redirects=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://{0}:81'.format(device_host)


def test_internal_web_open(device_host):

    response = requests.get('http://{0}:81'.format(device_host))
    assert response.status_code == 200


def test_activate_device(auth, device_host):

    email, password, domain, release = auth
    global LOGS_SSH_PASSWORD
    response = requests.post('http://{0}:81/rest/activate'.format(device_host),
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': 'user1', 'device_password': 'password1'})
    assert response.status_code == 200, response.text
    LOGS_SSH_PASSWORD = 'password1'


def test_reactivate(auth, device_host):
    email, password, domain, release = auth
    response = requests.post('http://{0}:81/rest/activate'.format(device_host),
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': DEVICE_USER, 'device_password': DEVICE_PASSWORD})
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = DEVICE_PASSWORD


def test_public_web_unauthorized_browser_redirect(device_host):
    response = requests.get('http://{0}/rest/user'.format(device_host), allow_redirects=False)
    assert response.status_code == 302


def test_public_web_unauthorized_ajax_not_redirect(device_host):
    response = requests.get('http://{0}/rest/user'.format(device_host),
                            allow_redirects=False, headers={'X-Requested-With': 'XMLHttpRequest'})
    assert response.status_code == 401


def test_running_platform_web(device_host):
    print(check_output('nc -zv -w 1 {0} 80'.format(device_host), shell=True))


def test_platform_rest(device_host):
    session = requests.session()
    session.mount('http://{0}'.format(device_host), HTTPAdapter(max_retries=5))
    response = session.get('http://{0}'.format(device_host), timeout=60)
    assert response.status_code == 200

def test_app_unix_socket(app_dir, device_host):
    run_scp('{0}/nginx.app.test.conf root@{1}:/'.format(DIR, device_host), throw=False, password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, '{0}/nginx/sbin/nginx -c /nginx.app.test.conf -g \'error_log {1}/log/nginx_app_error.log warn;\''.format(app_dir, DIR), password=DEVICE_PASSWORD)
    response = requests.get('http://unix_socket_app.{0}'.format(device_host), timeout=60)
    assert response.status_code == 200
    assert response.text == 'OK'


# def test_external_mode(auth, public_web_session, user_domain, device_host):
#
#     email, password, domain, release = auth
#
#     run_ssh(device_host, 'cp /integration/event/on_domain_change.py /opt/app/platform/bin', password=DEVICE_PASSWORD)
#
#     response = public_web_session.get('http://{0}/rest/settings/external_access'.format(device_host))
#     assert '"external_access": false' in response.text
#     assert response.status_code == 200
#
#     response = public_web_session.get('http://{0}/rest/settings/set_external_access'.format(device_host),
#                                       params={'external_access': 'true'})
#     assert '"success": true' in response.text
#     assert response.status_code == 200
#
#     response = public_web_session.get('http://{0}/rest/settings/external_access'.format(device_host))
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


def test_certbot_cli(app_dir, device_host):
    run_ssh(device_host, '{0}/bin/certbot --help'.format(app_dir), password=DEVICE_PASSWORD)
    run_ssh(device_host, '{0}/bin/certbot --help nginx'.format(app_dir), password=DEVICE_PASSWORD)


def test_openssl_cli(app_dir, device_host):
    run_ssh(device_host, '{0}/openssl/bin/openssl --help'.format(app_dir), password=DEVICE_PASSWORD)


def test_external_https_mode_with_certbot(public_web_session, device_host):

    response = public_web_session.get('http://{0}/rest/access/set_access'.format(device_host),
                                      params={'is_https': 'true', 'upnp_enabled': 'false', 'external_access': 'false', 'public_ip': 0, 'public_port': 0 })
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_show_https_certificate(device_host):
    run_ssh(device_host, "echo | "
            "openssl s_client -showcerts -servername localhost -connect localhost:443 2>/dev/null | "
            "openssl x509 -inform pem -noout -text", password=DEVICE_PASSWORD)


def test_access(public_web_session, device_host):
    response = public_web_session.get('http://{0}/rest/access/access'.format(device_host))
    print(response.text)
    assert '"success": true' in response.text
    assert '"upnp_enabled": false' in response.text
    assert response.status_code == 200


def test_network_interfaces(public_web_session, device_host):
    response = public_web_session.get('http://{0}/rest/access/network_interfaces'.format(device_host))
    print(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_device_url(public_web_session, device_host):
    response = public_web_session.get('http://{0}/rest/settings/device_url'.format(device_host))
    print(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_activate_url(public_web_session, device_host):
    response = public_web_session.get('http://{0}/rest/settings/activate_url'.format(device_host))
    print(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_hook_override(public_web_session, conf_dir, service_prefix, device_host):

    run_ssh(device_host, "sed -i 's#hooks_root.*#hooks_root: /integration#g' {0}/config/platform.cfg".format(conf_dir),
            password=DEVICE_PASSWORD)

    run_ssh(device_host, 'systemctl restart {0}platform.uwsgi-public'.format(service_prefix), password=DEVICE_PASSWORD)

    wait_for_rest(public_web_session, device_host, '/', 200)


def test_protocol(auth, public_web_session, device_host):

    email, password, domain, release = auth
 
    response = public_web_session.get('http://{0}/rest/access/access'.format(device_host))
    assert '"is_https": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://{0}/rest/access/set_access'.format(device_host),
                                      params={'is_https': 'true', 'upnp_enabled': 'false',
                                              'external_access': 'false', 'public_ip': 0, 'public_port': 0})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://{0}/rest/access/access'.format(device_host))
    assert '"is_https": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://{0}/rest/access/set_access'.format(device_host),
                                      params={'is_https': 'false', 'upnp_enabled': 'false',
                                              'external_access': 'false', 'public_ip': 0, 'public_port': 0})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://{0}/rest/access/access'.format(device_host))
    assert '"is_https": false' in response.text
    assert response.status_code == 200

    assert run_ssh(device_host, 'cat /tmp/on_domain_change.log',
                   password=DEVICE_PASSWORD) == '{0}.{1}'.format(domain, SYNCLOUD_INFO)


def test_cron_job(app_dir, ssh_env_vars, device_host):
    assert '"success": true' in run_ssh(device_host, '{0}/bin/insider sync_all'.format(app_dir),
                                        password=DEVICE_PASSWORD, env_vars=ssh_env_vars)


def test_installed_apps(public_web_session, device_host):
    response = public_web_session.get('http://{0}/rest/installed_apps'.format(device_host))
    assert response.status_code == 200


def test_do_not_cache_static_files_as_we_get_stale_ui_on_upgrades(public_web_session, device_host):

    response = public_web_session.get('http://{0}/settings.html'.format(device_host))
    cache_control = response.headers['Cache-Control']
    assert 'no-cache' in cache_control
    assert 'max-age=0' in cache_control


def test_installer_upgrade(public_web_session, device_host):
    __upgrade(public_web_session, 'sam', device_host)


@pytest.yield_fixture(scope='function')
def loop_device(device_host):
    dev_file = '/tmp/disk'
    loop_device_cleanup(device_host, dev_file, password=DEVICE_PASSWORD)

    print('adding loop device')
    run_ssh(device_host, 'dd if=/dev/zero bs=1M count=10 of={0}'.format(dev_file), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'sync', password=DEVICE_PASSWORD)
    run_ssh(device_host, 'ls -la {0}'.format(dev_file), password=DEVICE_PASSWORD)
    loop = run_ssh(device_host, 'losetup -f --show {0}'.format(dev_file), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'file -s {0}'.format(loop), password=DEVICE_PASSWORD)

    yield loop

    loop_device_cleanup(device_host, dev_file, password=DEVICE_PASSWORD)


def disk_writable(device_host):
    run_ssh(device_host, 'ls -la /data/', password=DEVICE_PASSWORD)
    run_ssh(device_host, "touch /data/platform/test.file", password=DEVICE_PASSWORD)


def test_udev_script(app_dir, device_host):
    run_ssh(device_host, '{0}/bin/check_external_disk'.format(app_dir), password=DEVICE_PASSWORD)


@pytest.mark.parametrize("fs_type", ['ext2', 'ext3', 'ext4'])
def test_public_settings_disk_add_remove(loop_device, public_web_session, fs_type, device_host):
    disk_create(loop_device, fs_type, device_host)
    assert disk_activate(loop_device,  public_web_session, device_host) == '/opt/disk/external/platform'
    disk_writable(device_host)
    assert disk_deactivate(loop_device, public_web_session, device_host) == '/opt/disk/internal/platform'


def test_disk_physical_remove(loop_device, public_web_session, device_host):
    disk_create(loop_device, 'ext4', device_host)
    assert disk_activate(loop_device,  public_web_session, device_host) == '/opt/disk/external/platform'
    loop_device_cleanup(device_host, '/opt/disk/external', password=DEVICE_PASSWORD)
    run_ssh(device_host, 'udevadm trigger --action=remove -y {0}'.format(loop_device.split('/')[2]),
            password=DEVICE_PASSWORD)
    run_ssh(device_host, 'udevadm settle', password=DEVICE_PASSWORD)
    assert current_disk_link(device_host) == '/opt/disk/internal/platform'


def disk_create(loop_device, fs, device_host):
    run_ssh(device_host, 'mkfs.{0} {1}'.format(fs, loop_device), password=DEVICE_PASSWORD)

    run_ssh(device_host, 'rm -rf /tmp/test', password=DEVICE_PASSWORD)
    run_ssh(device_host, 'mkdir /tmp/test', password=DEVICE_PASSWORD)

    run_ssh(device_host, 'mount {0} /tmp/test'.format(loop_device), password=DEVICE_PASSWORD)
    for mount in run_ssh(device_host, 'mount', debug=True, password=DEVICE_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh(device_host, 'umount {0}'.format(loop_device), password=DEVICE_PASSWORD)


def disk_activate(loop_device, public_web_session, device_host):

    response = public_web_session.get('http://{0}/rest/settings/disks'.format(device_host))
    print response.text
    assert loop_device in response.text
    assert response.status_code == 200

    response = public_web_session.get('http://{0}/rest/settings/disk_activate'.format(device_host),
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link(device_host)


def disk_deactivate(loop_device, public_web_session, device_host):
    response = public_web_session.get('http://{0}/rest/settings/disk_deactivate'.format(device_host),
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link(device_host)


def current_disk_link(device_host):
    return run_ssh(device_host, 'cat /tmp/on_disk_change.log', password=DEVICE_PASSWORD)


def test_internal_web_id(device_host):

    response = requests.get('http://{0}:81/rest/id'.format(device_host))
    assert 'mac_address' in response.text
    assert response.status_code == 200


def test_if_cron_is_enabled_after_install(device_host):
    cron_is_enabled_after_install(device_host)


def cron_is_enabled_after_install(device_host):
    crontab = run_ssh(device_host, "crontab -l", password=DEVICE_PASSWORD)
    assert len(crontab.splitlines()) == 1
    assert 'cron' in crontab, crontab
    assert not crontab.startswith('#'), crontab


def test_local_upgrade(app_archive_path, installer, device_host):
    if installer == 'sam':
        local_remove(device_host, DEVICE_PASSWORD, installer, 'platform')
        time.sleep(3)
        local_install(device_host, DEVICE_PASSWORD, app_archive_path, installer)
    else:
        local_install(device_host, DEVICE_PASSWORD, app_archive_path, installer)


def test_public_web_platform_upgrade(public_web_session, device_host):
    __upgrade(public_web_session, 'system', device_host)


def __upgrade(public_web_session, upgrade_type, device_host):

    public_web_session.get('http://{0}/rest/settings/{1}_upgrade'.format(device_host, upgrade_type))
    wait_for_sam(public_web_session, device_host)


def test_reinstall_local_after_upgrade(app_archive_path, installer, device_host):
    local_install(device_host, DEVICE_PASSWORD, app_archive_path, installer)


def test_if_cron_is_enabled_after_upgrade(device_host):
    cron_is_enabled_after_install(device_host)


def test_nginx_performance(device_host):
    print(check_output('ab -c 1 -n 1000 http://{0}/ping'.format(device_host), shell=True))


def test_nginx_plus_flask_performance(device_host):
    print(check_output('ab -c 1 -n 1000 http://{0}:81/rest/id'.format(device_host), shell=True))
