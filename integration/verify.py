import json
from os.path import dirname, join
from subprocess import check_output

import pytest
import requests
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.exceptions import InsecureRequestWarning
from syncloudlib.http import wait_for_response
from syncloudlib.integration.hosts import add_host_alias_by_ip
from syncloudlib.integration.installer import local_install, wait_for_installer
from syncloudlib.integration.loop import loop_device_cleanup
from syncloudlib.integration.ssh import run_ssh

requests.packages.urllib3.disable_warnings(InsecureRequestWarning)
from syncloudlib.application.storage import init_storage
from syncloudlib.application.service import restart
from syncloudlib.application.paths import get_data_dir, get_app_dir
from syncloudlib.application.config import set_dkim_key, get_dkim_key
from syncloudlib.application.urls import get_app_url

DIR = dirname(__file__)
TMP_DIR = '/tmp/syncloud'
DEFAULT_LOGS_SSH_PASSWORD = "syncloud"
LOGS_SSH_PASSWORD = DEFAULT_LOGS_SSH_PASSWORD


@pytest.fixture(scope="session")
def app_data_dir():
    return '/var/snap/{0}/common'.format('app')


@pytest.fixture(scope="session")
def module_setup(request, data_dir, device, app_dir, artifact_dir):
    def module_teardown():
        device.scp_from_device('{0}/config'.format(data_dir), artifact_dir)
        device.scp_from_device('{0}/config.runtime'.format(data_dir), artifact_dir)
        device.run_ssh('mkdir {0}'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('snap run platform.cli ipv4 public > {0}/cli.ipv4.public.log'.format(TMP_DIR), throw=False)
        device.run_ssh('snap run platform.cli config list > {0}/cli.config.list.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ps auxfw > {0}/ps.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la {0}/ > {1}/app.data.ls.log'.format(data_dir, TMP_DIR), throw=False)
        device.run_ssh('ls -la {0}/ > {1}/app.ls.log'.format(app_dir, TMP_DIR), throw=False)
        device.run_ssh('ls -la {0}/www > {1}/app.www.ls.log'.format(app_dir, TMP_DIR), throw=False)
        device.run_ssh('ls -la /data/platform/backup > {0}/data.platform.backup.ls.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), artifact_dir)
        device.scp_from_device('{0}/log/*'.format(data_dir), artifact_dir)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, device, app, domain, device_host):
    device.run_ssh('date', retries=100, throw=True)
    device.scp_to_device(DIR, '/', throw=True)
    device.run_ssh('mkdir /etc/syncloud', throw=True)
    device.scp_to_device(join(DIR, 'id.cfg'), '/etc/syncloud', throw=True)
    device.run_ssh('mkdir /log', throw=True)
    device.run_ssh('snap remove platform', throw=True)
    add_host_alias_by_ip(app, domain, device_host)
    add_host_alias_by_ip("app", domain, device_host)


def test_install(app_archive_path, device_host):
    local_install(device_host, DEFAULT_LOGS_SSH_PASSWORD, app_archive_path)


def test_install_testapp(device_host):
    local_install(device_host, DEFAULT_LOGS_SSH_PASSWORD, join(DIR, 'testapp', 'testapp.snap'))


def test_cryptography_openssl_version(device_host):
    run_ssh(device_host, "/snap/platform/current/python/bin/python "
                         "-c 'from cryptography.hazmat.backends.openssl.backend "
                         "import backend; print(backend.openssl_version_text())'", password=LOGS_SSH_PASSWORD)


def test_http_port_validation_url(device_host):
    response = requests.get('http://{0}/ping'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert response.text == 'OK'


def test_https_port_validation_url(device_host):
    response = requests.get('https://{0}/ping'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert response.text == 'OK'


def test_non_activated_device_login_redirect_to_activation(device_host):
    response = requests.post('https://{0}/rest/login'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 501


def test_activation_status_false(device_host):
    response = requests.get('https://{0}/rest/activation_status'.format(device_host),
                            allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert not json.loads(response.text)["activated"], response.text


def test_id_redirect_backward_compatibility(device_host):
    response = requests.get('http://{0}:81/rest/id'.format(device_host), allow_redirects=False)
    assert response.status_code == 200, response.text
    response_json = json.loads(response.text)
    assert 'data' in response_json
    assert 'success' in response_json
    assert 'mac_address' in response_json['data']
    assert 'title' in response_json['data']
    assert 'name' in response_json['data']


def test_id_before_activation(device_host):
    response = requests.get('https://{0}/rest/id'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 200, response.text
    response_json = json.loads(response.text)
    assert 'data' in response_json
    assert 'success' in response_json
    assert 'mac_address' in response_json['data']
    assert 'title' in response_json['data']
    assert 'name' in response_json['data']


def test_set_redirect(device, main_domain):
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))


def test_activate_custom(device, device_host, main_domain):
    response = requests.post('https://{0}/rest/activate/custom'.format(device_host),
                             json={'domain': 'example.com',
                                   'device_username': 'user1',
                                   'device_password': DEFAULT_LOGS_SSH_PASSWORD}, verify=False)
    assert response.status_code == 200, response.text
    device.run_ssh('rm /var/snap/platform/common/platform.db')
    device.run_ssh('ls -la /var/snap/platform/common')
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))


def test_activate_premium(device, device_host, main_domain, redirect_user, redirect_password, arch):
    response = requests.post('https://{0}/rest/activate/premium'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': '{}-syncloudexample.com'.format(arch),
                                   'device_username': 'user1',
                                   'device_password': DEFAULT_LOGS_SSH_PASSWORD}, verify=False)
    assert response.status_code == 200, response.text
    device.run_ssh('rm /var/snap/platform/common/platform.db')
    device.run_ssh('ls -la /var/snap/platform/common')
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))


def test_activate_device(device_host, domain, redirect_user, redirect_password):
    response = requests.post('https://{0}/rest/activate/free'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': domain,
                                   'device_username': 'user1',
                                   'device_password': DEFAULT_LOGS_SSH_PASSWORD}, verify=False)
    assert response.status_code == 200, response.text


def test_reactivate_activated_device(device_host, domain, device_user, device_password,
                                     redirect_user, redirect_password):
    response = requests.post('https://{0}/rest/activate/free'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': domain,
                                   'device_username': device_user,
                                   'device_password': device_password}, allow_redirects=False, verify=False)
    assert response.status_code == 502, response.text


def test_drop_activation(device, main_domain):
    device.run_ssh('rm /var/snap/platform/common/platform.db')
    device.run_ssh('ls -la /var/snap/platform/common')
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))


def test_reactivate_good(device_host, domain, device_user, device_password,
                         redirect_user, redirect_password, device):
    response = requests.post('https://{0}/rest/activate/free'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': domain,
                                   'device_username': device_user,
                                   'device_password': device_password}, verify=False)
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = device_password
    device.ssh_password = device_password


def test_deactivate(device, device_host):
    response = device.login().post('https://{0}/rest/settings/deactivate'.format(device_host), verify=False,
                                   allow_redirects=False)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_activation_status_false_after_deactivate(device_host):
    response = requests.get('https://{0}/rest/activation_status'.format(device_host), allow_redirects=False,
                            verify=False)
    assert response.status_code == 200
    assert not json.loads(response.text)["activated"], response.text


def test_redirect_info(device_host, main_domain):
    response = requests.get('https://{0}/rest/redirect_info'.format(device_host), allow_redirects=False,
                            verify=False)
    assert response.status_code == 200
    assert json.loads(response.text)["domain"] == main_domain, response.text


def test_reactivate_after_deactivate(device_host, domain, device_user, device_password,
                                     redirect_user, redirect_password, device):
    response = requests.post('https://{0}/rest/activate/free'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': domain,
                                   'device_username': device_user,
                                   'device_password': device_password}, verify=False)
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = device_password
    device.ssh_password = device_password


def test_activation_status_true(device_host):
    response = requests.get('https://{0}/rest/activation_status'.format(device_host), allow_redirects=False,
                            verify=False)
    assert response.status_code == 200
    assert json.loads(response.text)["activated"], response.text


def test_unauthorized(device_host):
    response = requests.get('https://{0}/rest/user'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 401


def test_running_platform_web(device_host):
    print(check_output('nc -zv -w 1 {0} 80'.format(device_host), shell=True).decode())


def test_platform_rest(device_host):
    session = requests.session()
    session.mount('https://{0}'.format(device_host), HTTPAdapter(max_retries=5))
    response = session.get('https://{0}'.format(device_host), timeout=60, verify=False)
    assert response.status_code == 200


def test_api_service_restart():
    status = restart('platform.nginx-public')
    assert 'OK' in status, status


def test_api_install_path(app_dir):
    response = get_app_dir('platform')
    assert app_dir in response, response


def test_api_config_dkim_key():
    response = set_dkim_key('dkim123')
    assert 'OK' in response, response

    response = get_dkim_key()
    assert 'dkim123' in response, response


def test_api_data_path(data_dir):
    response = get_data_dir('platform')
    assert data_dir in response, response


def test_api_url(app_domain):
    response = get_app_url('platform')
    assert app_domain in response, response


def test_certbot_cli(app_dir, device_host):
    output = run_ssh(device_host, '{0}/bin/certbot --help'.format(app_dir), password=LOGS_SSH_PASSWORD)
    assert not output.strip() == ""
    run_ssh(device_host, '{0}/bin/certbot --help nginx'.format(app_dir), password=LOGS_SSH_PASSWORD)


def test_openssl_cli(app_dir, device_host):
    run_ssh(device_host, '{0}/openssl/bin/openssl help'.format(app_dir), password=LOGS_SSH_PASSWORD)


def test_set_access_mode_with_certbot(device, device_host):
    response = device.login().post('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                   json={'upnp_enabled': False,
                                         'external_access': False,
                                         'public_ip': 0,
                                         'certificate_port': 0, 'access_port': 0})
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_testapp_access_change(device_host, domain):
    output = run_ssh(device_host, 'cat /var/snap/testapp/common/on_access_change', password=LOGS_SSH_PASSWORD)
    assert not output.strip() == "https://testapp.{0}".format(domain)


def test_testapp_access_change_hook(device_host):
    run_ssh(device_host, 'snap run testapp.access-change', password=LOGS_SSH_PASSWORD)


def test_show_https_certificate(device_host):
    run_ssh(device_host, "echo | "
                         "openssl s_client -showcerts -servername localhost -connect localhost:443 2>/dev/null | "
                         "openssl x509 -inform pem -noout -text", password=LOGS_SSH_PASSWORD)


def test_get_access(device, device_host):
    response = device.login().get('https://{0}/rest/access/access'.format(device_host), verify=False)
    print(response.text)
    assert '"success": true' in response.text
    assert '"upnp_enabled": false' in response.text
    assert response.status_code == 200


def test_network_interfaces(device, device_host):
    response = device.login().get('https://{0}/rest/access/network_interfaces'.format(device_host), verify=False)
    print(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_send_logs(device, device_host):
    response = device.login().post('https://{0}/rest/send_log?include_support=false'.format(device_host), verify=False)
    print(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_available_apps(device, device_host, artifact_dir):
    response = device.login().get('https://{0}/rest/available_apps'.format(device_host), verify=False)
    with open('{0}/rest.available_apps.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200
    assert len(json.loads(response.text)['apps']) > 1


def test_device_url(device, device_host, artifact_dir):
    response = device.login().get('https://{0}/rest/settings/device_url'.format(device_host), verify=False)
    with open('{0}/rest.settings.device_url.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_api_url_443(device, device_host, app_domain):
    response = device.login().get('https://{0}/rest/access/access'.format(device_host), verify=False)
    assert response.status_code == 200

    response = device.login().post('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                   json={'upnp_enabled': False,
                                         'external_access': False, 'public_ip': 0,
                                         'certificate_port': 443, 'access_port': 443})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = device.login().get('https://{0}/rest/access/access'.format(device_host), verify=False)
    assert response.status_code == 200
    url = get_app_url('platform')

    assert app_domain in url, url
    assert 'https' in url, url


def test_api_url_10000(device, device_host, app_domain):
    response = device.login().post('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                   json={'upnp_enabled': False,
                                         'external_access': False, 'public_ip': 0,
                                         'certificate_port': 80, 'access_port': 10000})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = device.login().get('https://{0}/rest/access/access'.format(device_host), verify=False)
    assert response.status_code == 200

    url = get_app_url('platform')

    assert app_domain in url, url
    assert 'https' in url, url


def test_set_access_error(device, device_host):
    response = device.login().post('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                   json={'upnp_enabled': False,
                                         'external_access': True, 'public_ip': 0,
                                         'certificate_port': 0, 'access_port': 0})
    assert '"success": false' in response.text
    assert 'Unable to verify open ports' in response.text
    assert response.status_code == 500, response.text


def test_cron(device):
    device.run_ssh('snap run platform.cli cron')


def test_real_certificate(app_dir, ssh_env_vars, device_host):
    run_ssh(device_host, '{0}/bin/generate_real_certificate.py'.format(app_dir),
            password=LOGS_SSH_PASSWORD, env_vars=ssh_env_vars)


def test_install_app(device, device_host):
    session = device.login()
    session.post('https://{0}/rest/install'.format(device_host), json={'app_id': 'files'}, verify=False)
    wait_for_installer(session, device_host)


def test_rest_installed_apps(device, device_host, artifact_dir):
    response = device.login().get('https://{0}/rest/installed_apps'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.installed_apps.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200
    assert len(json.loads(response.text)['apps']) == 2


def test_rest_installed_app(device, device_host, artifact_dir):
    response = device.login().get('https://{0}/rest/app?app_id=files'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.app.installed.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200


def test_rest_not_installed_app(device, device_host, artifact_dir):
    response = device.login().get('https://{0}/rest/app?app_id=nextcloud'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.app.not.installed.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200


def test_installer_upgrade(device, device_host):
    session = device.login()
    response = session.post('https://{0}/rest/installer/upgrade'.format(device_host), verify=False)
    assert response.status_code == 200, response.text
    wait_for_response(session, 'https://{0}/rest/job/status'.format(device_host),
                      lambda r: json.loads(r.text)['data'] == 'JobStatusIdle')

    response = session.post('https://{0}/rest/installer/upgrade'.format(device_host), verify=False)
    assert response.status_code == 200, response.text
    wait_for_response(session, 'https://{0}/rest/job/status'.format(device_host),
                      lambda r: json.loads(r.text)['data'] == 'JobStatusIdle')


def test_backup_app(device, artifact_dir, device_host):
    session = device.login()
    response = session.post('https://{0}/rest/backup/create'.format(device_host), json={'app': 'files'}, verify=False)
    assert response.status_code == 200
    assert json.loads(response.text)['success']

    wait_for_response(session, 'https://{0}/rest/job/status'.format(device_host),
                      lambda r: json.loads(r.text)['data'] == 'JobStatusIdle')

    response = device.http_get('/rest/backup/list')
    assert response.status_code == 200
    open('{0}/rest.backup.list.json'.format(artifact_dir), 'w').write(response.text)
    print(response.text)
    backup = json.loads(response.text)['data'][0]
    device.run_ssh('tar tvf {0}/{1}'.format(backup['path'], backup['file']))

    response = session.post(
        'https://{0}/rest/backup/restore'.format(device_host),
        json={'app': 'files', 'file': '{0}/{1}'.format(backup['path'], backup['file'])},
        verify=False)
    assert response.status_code == 200
    wait_for_response(session, 'https://{0}/rest/job/status'.format(device_host),
                      lambda r: json.loads(r.text)['data'] == 'JobStatusIdle')


def test_rest_backup_list(device, device_host, artifact_dir):
    response = device.login().get('https://{0}/rest/backup/list'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.backup.list.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert json.loads(response.text)['success']


@pytest.fixture(scope='function')
def loop_device(device_host):
    dev_file = '/tmp/disk'
    loop_device_cleanup(device_host, dev_file, password=LOGS_SSH_PASSWORD)

    print('adding loop device')
    run_ssh(device_host, 'dd if=/dev/zero bs=1M count=10 of={0}'.format(dev_file), password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'sync', password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'ls -la {0}'.format(dev_file), password=LOGS_SSH_PASSWORD)
    loop = run_ssh(device_host, 'losetup -f --show {0}'.format(dev_file), password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'losetup', password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'losetup -j {0} | grep {0}'.format(dev_file), password=LOGS_SSH_PASSWORD, retries=3)
    # run_ssh(device_host, 'file -s {0}'.format(loop), password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'sync', password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'partprobe {0}'.format(loop), password=LOGS_SSH_PASSWORD, retries=3)
    yield loop

    loop_device_cleanup(device_host, dev_file, password=LOGS_SSH_PASSWORD)


def disk_writable(device_host):
    run_ssh(device_host, 'ls -la /data/', password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, "touch /data/platform/test.file", password=LOGS_SSH_PASSWORD)


@pytest.mark.parametrize("fs_type", ['ext4'])
def test_public_settings_disk_add_remove(loop_device, device, fs_type, device_host):
    disk_create(loop_device, fs_type, device_host)
    assert disk_activate(loop_device, device, device_host) == '/opt/disk/external/platform'
    disk_writable(device_host)
    assert disk_deactivate(loop_device, device, device_host) == '/opt/disk/internal/platform'


def disk_create(loop, fs, device_host):
    tmp_disk = '/tmp/test'
    run_ssh(device_host, 'mkfs.{0} {1}'.format(fs, loop), password=LOGS_SSH_PASSWORD, retries=3)

    run_ssh(device_host, 'rm -rf {0}'.format(tmp_disk), password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'mkdir {0}'.format(tmp_disk), password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'sync', password=LOGS_SSH_PASSWORD)

    run_ssh(device_host, 'mount {0} {1}'.format(loop, tmp_disk), password=LOGS_SSH_PASSWORD, retries=3)
    for mount in run_ssh(device_host, 'mount', debug=True, password=LOGS_SSH_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh(device_host, 'umount {0}'.format(loop), password=LOGS_SSH_PASSWORD)


def disk_activate(loop, device, device_host):
    response = device.login().get('http://{0}/rest/settings/disks'.format(device_host))
    print(response.text)
    assert loop in response.text
    assert response.status_code == 200

    response = device.login().post('https://{0}/rest/settings/disk_activate'.format(device_host), verify=False,
                                   json={'device': loop})
    assert response.status_code == 200
    return current_disk_link(device)


def disk_deactivate(loop, device, device_host):
    response = device.login().post('https://{0}/rest/settings/disk_deactivate'.format(device_host), verify=False,
                                   json={'device': loop})
    assert response.status_code == 200
    return current_disk_link(device)


def current_disk_link(device):
    link = init_storage('platform', 'root')
    return device.run_ssh("realpath {}".format(link))


def test_if_cron_is_empty_after_install(device_host):
    cron_is_empty_after_install(device_host)


def cron_is_empty_after_install(device_host):
    crontab = run_ssh(device_host, "crontab -l", password=LOGS_SSH_PASSWORD, throw=False)
    assert 'no crontab for root' in crontab, crontab


def test_settings_versions(device_host, device, artifact_dir):
    response = device.login().get('https://{0}/rest/settings/versions'.format(device_host), allow_redirects=False,
                                  verify=False)
    with open('{0}/rest.settings.versions.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)

    assert response.status_code == 200, response.text


def test_local_upgrade(app_archive_path, device_host):
    local_install(device_host, LOGS_SSH_PASSWORD, app_archive_path)


def test_reinstall_local_after_upgrade(app_archive_path, device_host):
    local_install(device_host, LOGS_SSH_PASSWORD, app_archive_path)


def test_if_cron_is_empty_after_upgrade(device_host):
    cron_is_empty_after_install(device_host)


def test_nginx_performance(device_host):
    print(check_output('ab -c 1 -n 1000 https://{0}/ping'.format(device_host), shell=True).decode())


def test_nginx_plus_flask_performance(device_host):
    print(check_output('ab -c 1 -n 1000 https://{0}/rest/id'.format(device_host), shell=True).decode())

