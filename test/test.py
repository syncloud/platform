import re
import json
import time
from os.path import dirname, join
from subprocess import check_output

import pytest
import requests
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.exceptions import InsecureRequestWarning
from syncloudlib.http import wait_for_response
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.installer import local_install
from syncloudlib.integration.loop import loop_device_cleanup
from syncloudlib.integration.ssh import run_ssh

requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

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
        device.run_ssh('journalctl > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('snap run platform.cli ipv4 public > {0}/cli.ipv4.public.log'.format(TMP_DIR), throw=False)
        device.run_ssh('snap run platform.cli config list > {0}/cli.config.list.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ps auxfw > {0}/ps.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/platform/current/config/authelia/config.yml {0}/authelia.config.yml.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap/platform/common/ > {0}/snap.common.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap/platform/current/ > {0}/snap.data.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /snap/platform/current/ > {0}/snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la {0}/www > {1}/app.www.ls.log'.format(app_dir, TMP_DIR), throw=False)
        device.run_ssh('ls -la /data/platform/backup > {0}/data.platform.backup.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('df -h > {0}/df.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), artifact_dir)
        device.scp_from_device('{0}/log/*'.format(data_dir), artifact_dir)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, device, app, domain, device_host, full_domain):
    add_host_alias(app, device_host, full_domain)
    add_host_alias(app, device_host, domain)
    add_host_alias("app", device_host, domain)
    add_host_alias("testapp", device_host, full_domain)
    device.run_ssh("echo '127.0.0.1 auth.{0}' >> /etc/hosts".format(full_domain))
    device.run_ssh('mkdir {0}'.format(TMP_DIR), throw=False)
    device.run_ssh('date', retries=100, throw=True)
    device.scp_to_device(DIR, '/', throw=True)
    device.run_ssh('/test/install-snapd.sh', throw=True)
    device.run_ssh('mkdir /etc/syncloud', throw=True)
    device.run_ssh('rm -rf /usr/lib/sasl2', throw=True)
    device.scp_to_device(join(DIR, 'id.cfg'), '/etc/syncloud', throw=True)
    device.run_ssh('mkdir /log', throw=True)
    

def test_install(app_archive_path, device_host):
    local_install(device_host, DEFAULT_LOGS_SSH_PASSWORD, app_archive_path)


def test_https_port_validation_url(device_host):
    response = requests.get('https://{0}/ping'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert response.text == 'OK'


def test_non_activated_device_login_redirect_to_activation(full_domain):
    def login():
        response = requests.get('https://{0}/rest/oidc/login'.format(full_domain), verify=False)
        if response.status_code != 501:
            raise Exception()
    retry(login)


def test_activation_status_false(device_host):
    response = requests.get('https://{0}/rest/activation/status'.format(device_host),
                            allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert not json.loads(response.text)["data"], response.text


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


def test_activate_custom(device, device_host, main_domain):
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))
    device.run_ssh('snap run platform.cli config set redirect.api_url http://api.redirect')
    device.run_ssh('snap run platform.cli config set certbot.staging true')
    response = requests.post('https://{0}/rest/activate/custom'.format(device_host),
                             json={'domain': 'example.com',
                                   'device_username': 'user1',
                                   'device_password': DEFAULT_LOGS_SSH_PASSWORD}, verify=False)
    assert response.status_code == 200, response.text
    

def test_activate_premium(device, device_host, main_domain, redirect_user, redirect_password, arch):
    device.run_ssh('rm /var/snap/platform/current/platform.db')
    device.run_ssh('ls -la /var/snap/platform/current')
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))
    device.run_ssh('snap run platform.cli config set redirect.api_url http://api.redirect')
    device.run_ssh('snap run platform.cli config set certbot.staging true')
    response = requests.post('https://{0}/rest/activate/managed'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': '{}-syncloudexample.com'.format(arch),
                                   'device_username': 'user1',
                                   'device_password': DEFAULT_LOGS_SSH_PASSWORD}, verify=False)
    assert response.status_code == 200, response.text
    

def test_activate_device(device, device_host, main_domain, full_domain, redirect_user, redirect_password):
    device.run_ssh('rm /var/snap/platform/current/platform.db')
    device.run_ssh('ls -la /var/snap/platform/current')
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))
    device.run_ssh('snap run platform.cli config set redirect.api_url http://api.redirect')
    device.run_ssh('snap run platform.cli config set certbot.staging true')
    response = requests.post('https://{0}/rest/activate/managed'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': full_domain,
                                   'device_username': 'user1',
                                   'device_password': DEFAULT_LOGS_SSH_PASSWORD}, verify=False)
    assert response.status_code == 200, response.text


def test_reactivate_activated_device(device_host, full_domain, device_user, device_password,
                                     redirect_user, redirect_password):
    response = requests.post('https://{0}/rest/activate/managed'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': full_domain,
                                   'device_username': device_user,
                                   'device_password': device_password}, allow_redirects=False, verify=False)
    assert response.status_code == 502, response.text


def test_drop_activation(device, main_domain):
    device.run_ssh('rm /var/snap/platform/current/platform.db')
    device.run_ssh('ls -la /var/snap/platform/current')
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))
    device.run_ssh('snap run platform.cli config set redirect.api_url http://api.redirect')
    device.run_ssh('snap run platform.cli config set certbot.staging true')


def test_reactivate_good(device_host, full_domain, device_user, device_password,
                         redirect_user, redirect_password, device):
    response = requests.post('https://{0}/rest/activate/managed'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': full_domain,
                                   'device_username': device_user,
                                   'device_password': device_password}, verify=False)
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = device_password
    device.ssh_password = device_password


def wait_for_activation(domain):
    def check():
        response = requests.get('https://{0}/rest/activation/status'.format(domain),
                                allow_redirects=False, verify=False)
        if response.status_code != 200:
            raise Exception()
        if not json.loads(response.text)["data"]:
            raise Exception()
    retry(check)


def test_deactivate(device, domain):
    wait_for_activation(domain)
    response = device.login_v2().post('https://{0}/rest/deactivate'.format(domain), verify=False,
                                   allow_redirects=False)
    assert '"success":true' in response.text
    assert response.status_code == 200


def test_activation_status_false_after_deactivate(device_host):
    response = requests.get('https://{0}/rest/activation/status'.format(device_host), allow_redirects=False,
                            verify=False)
    assert response.status_code == 200
    assert not json.loads(response.text)["data"], response.text


def test_redirect_info(device_host, main_domain):
    response = requests.get('https://{0}/rest/redirect_info'.format(device_host), allow_redirects=False,
                            verify=False)
    assert response.status_code == 200
    assert json.loads(response.text)['data']["domain"] == main_domain, response.text


def test_reactivate_after_deactivate(device_host, full_domain, device_user, device_password,
                                     redirect_user, redirect_password, device):
    response = requests.post('https://{0}/rest/activate/managed'.format(device_host),
                             json={'redirect_email': redirect_user,
                                   'redirect_password': redirect_password,
                                   'domain': full_domain,
                                   'device_username': device_user,
                                   'device_password': device_password}, verify=False)
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = device_password
    device.ssh_password = device_password


def test_activation_status_true(device_host):
    response = requests.get('https://{0}/rest/activation/status'.format(device_host), allow_redirects=False,
                            verify=False)
    assert response.status_code == 200
    assert json.loads(response.text)["data"], response.text


def test_install_ca_cert(device_host, full_domain):
    run_ssh(device_host, 'cp /var/snap/platform/current/syncloud.ca.crt /usr/local/share/ca-certificates/syncloud.crt', password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'update-ca-certificates', password=LOGS_SSH_PASSWORD)
    run_ssh(device_host, 'curl https://auth.{0}'.format(full_domain), password=LOGS_SSH_PASSWORD)


def test_install_testapp(device_host):
    local_install(device_host, LOGS_SSH_PASSWORD, join(DIR, 'testapp', 'testapp.snap'))


def test_unauthorized(device_host):
    response = requests.get('https://{0}/rest/user'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 401


def test_running_platform_web(device_host):
    print(check_output('nc -zv -w 1 {0} 443'.format(device_host), shell=True).decode())


def test_platform_rest(device_host):
    session = requests.session()
    session.mount('https://{0}'.format(device_host), HTTPAdapter(max_retries=5))
    response = session.get('https://{0}'.format(device_host), timeout=60, verify=False)
    assert response.status_code == 200


def test_api(device):
    time.sleep(10) # start-limit-hit
    device.scp_to_device(join(DIR, "api/api.test"), '/', throw=True)
    device.run_ssh('/api.test')


def test_oidc_discovery(full_domain):
    response = requests.get('https://auth.{0}/.well-known/openid-configuration'.format(full_domain), verify=False)
    assert response.status_code == 200, response.text
    config = response.json()
    assert 'authorization_endpoint' in config
    assert 'token_endpoint' in config


def test_oauth2_jwks(full_domain):
    discovery = requests.get('https://auth.{0}/.well-known/openid-configuration'.format(full_domain), verify=False).json()
    jwks_uri = discovery['jwks_uri']
    response = requests.get(jwks_uri, verify=False)
    assert response.status_code == 200, response.text
    jwks = response.json()
    assert 'keys' in jwks


def test_custom_proxy(device, device_host, full_domain):
    device.run_ssh('nohup /test/externalapp/externalapp > /tmp/syncloud/externalapp.log 2>&1 &', throw=False)
    time.sleep(2)
    add_host_alias("externalapp", device_host, full_domain)
    session = device.login_v2()
    response = session.post('https://{0}/rest/proxy_custom/add'.format(device_host),
                            json={'name': 'externalapp', 'host': 'localhost', 'port': 8585},
                            verify=False)
    assert response.status_code == 200, response.text
    assert json.loads(response.text)["success"], response.text

    response = session.get('https://{0}/rest/proxy_custom/list'.format(device_host), verify=False)
    assert response.status_code == 200, response.text
    proxies = json.loads(response.text)["data"]
    assert len(proxies) == 1
    assert proxies[0]["name"] == "externalapp"

    def check_proxy():
        response = requests.get('https://externalapp.{0}'.format(full_domain), verify=False)
        assert response.status_code == 200, response.text
        assert response.text == "external", response.text
    retry(check_proxy)

    response = session.post('https://{0}/rest/proxy_custom/remove'.format(device_host),
                            json={'name': 'externalapp'},
                            verify=False)
    assert response.status_code == 200, response.text


def test_custom_proxy_cli(device):
    device.run_ssh('snap run platform.cli proxy add --name cliproxy --host localhost --port 9090 --https')
    output = device.run_ssh('snap run platform.cli proxy list')
    proxies = json.loads(output)
    assert len(proxies) == 1
    assert proxies[0]["name"] == "cliproxy"
    assert proxies[0]["host"] == "localhost"
    assert proxies[0]["port"] == 9090
    assert proxies[0]["https"] is True

    device.run_ssh('snap run platform.cli proxy remove --name cliproxy')
    output = device.run_ssh('snap run platform.cli proxy list')
    proxies = json.loads(output)
    assert len(proxies) == 0


def test_testapp_access_change(device_host, full_domain):
    output = run_ssh(device_host, 'cat /var/snap/testapp/current/config/authelia-location.conf', password=LOGS_SSH_PASSWORD)
    assert 'unix:/var/snap/platform/current/authelia.socket:' in output, output


def test_testapp_access_change_hook(device_host):
    run_ssh(device_host, 'snap run testapp.access-change', password=LOGS_SSH_PASSWORD)


def test_testapp_session_unauthorized(full_domain):
    response = requests.get(
        'https://testapp.{0}/'.format(full_domain),
        verify=False,
        allow_redirects=False)
    assert response.status_code == 302, "expected redirect to login, got {0}".format(response.status_code)
    assert 'auth.{0}'.format(full_domain) in response.headers.get('Location', ''), response.headers.get('Location', '')


def test_testapp_session_authorized(full_domain, device_user, device_password):
    session = requests.session()
    session.post(
        'https://auth.{0}/api/firstfactor'.format(full_domain),
        json={'username': device_user, 'password': device_password},
        verify=False)
    response = session.get(
        'https://testapp.{0}/'.format(full_domain),
        verify=False,
        allow_redirects=False)
    assert response.status_code == 200, "expected 200 with session, got {0}: {1}".format(response.status_code, response.text)
    assert 'session protected page' in response.text, response.text


def test_testapp_basic_unauthorized(full_domain):
    response = requests.get(
        'https://testapp.{0}/basic'.format(full_domain),
        verify=False,
        allow_redirects=False)
    assert response.status_code == 401, "expected 401 without credentials, got {0}".format(response.status_code)


def test_testapp_basic_authorized(full_domain, device_user, device_password):
    response = requests.get(
        'https://{0}:{1}@testapp.{2}/basic'.format(device_user, device_password, full_domain),
        verify=False,
        allow_redirects=False)
    assert response.status_code == 200, "expected 200 with basic auth, got {0}: {1}".format(response.status_code, response.text)
    assert 'session protected page' in response.text, response.text


def test_get_access(device, domain):
    response = device.login_v2().get('https://{0}/rest/access'.format(domain), verify=False)
    print(response.text)
    assert json.loads(response.text)["success"]
    assert json.loads(response.text)["data"]["ipv4_enabled"]
    assert response.status_code == 200


def test_installer_status(device, device_host):
    session = device.login_v2()
    for attempt in range(30):
        response = session.get('https://{0}/rest/installer/status'.format(device_host), allow_redirects=False, verify=False)
        assert response.status_code == 200
        if not json.loads(response.text)["data"]["is_running"]:
            return
        print('installer still running (attempt {0}/30): {1}'.format(attempt + 1, response.text))
        time.sleep(2)
    assert False, 'installer still running after 60s: {0}'.format(response.text)


def test_network_interfaces(device, domain):
    response = device.login_v2().get('https://{0}/rest/network/interfaces'.format(domain), verify=False)
    print(response.text)
    assert json.loads(response.text)["success"]
    assert response.status_code == 200


def test_send_logs(device, domain):
    response = device.login_v2().post('https://{0}/rest/logs/send?include_support=false'.format(domain), verify=False)
    print(response.text)
    assert json.loads(response.text)["success"]
    assert response.status_code == 200


def test_proxy_image(device, domain):
    response = device.login_v2().get('https://{0}/rest/proxy/image?channel=stable&app=files'.format(domain), verify=False)
    assert response.status_code == 200


def test_available_apps(device, domain, artifact_dir):
    response = device.login_v2().get('https://{0}/rest/apps/available'.format(domain), verify=False)
    with open('{0}/rest.available_apps.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200
    assert len(json.loads(response.text)['data']) > 0


def test_device_url(device, domain, artifact_dir, full_domain):
    response = device.login_v2().get('https://{0}/rest/device/url'.format(domain), verify=False)
    with open('{0}/rest.device.url.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert json.loads(response.text)["success"]
    assert response.status_code == 200
    assert json.loads(response.text)["data"] == 'https://{}'.format(full_domain), response.text


def test_api_url_443(device, domain):
    response = device.login_v2().get('https://{0}/rest/access'.format(domain), verify=False)
    assert response.status_code == 200

    response = device.login_v2().post('https://{0}/rest/access'.format(domain), verify=False,
                                   json={'ipv4_enabled': False,
                                         'ipv4_public': False,
                                         'access_port': 443})
    assert json.loads(response.text)["success"]
    assert response.status_code == 200

    response = device.login_v2().get('https://{0}/rest/access'.format(domain), verify=False)
    assert response.status_code == 200


def test_api_url_10000(device, domain):
    response = device.login_v2().post('https://{0}/rest/access'.format(domain), verify=False,
                                   json={'ipv4_enabled': False,
                                         'ipv4_public': False,
                                         'access_port': 10000})
    assert json.loads(response.text)["success"]
    assert response.status_code == 200

    response = device.login_v2().get('https://{0}/rest/access'.format(domain), verify=False)
    assert response.status_code == 200

    response = device.login_v2().post('https://{0}/rest/access'.format(domain), verify=False,
                                   json={'ipv4_enabled': False,
                                         'ipv4_public': False,
                                         'access_port': 443})
    assert json.loads(response.text)["success"]


def test_cron(device):
    device.run_ssh('snap run platform.cli cron')


def test_rest_installed_apps(device, domain, artifact_dir):
    response = device.login_v2().get('https://{0}/rest/apps/installed'.format(domain), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.installed_apps.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200
    assert len(json.loads(response.text)['data']) == 1


def test_rest_not_installed_app(device, domain, artifact_dir):
    response = device.login_v2().get('https://{0}/rest/app?app_id=files'.format(domain), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.app.not.installed.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200


def test_rest_platform_version(device, domain, artifact_dir):
    response = device.login_v2().get('https://{0}/rest/app?app_id=platform'.format(domain), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.platform.version.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200



def wait_for_jobs(domain, session):
    wait_for_response(session, 'https://{0}/rest/job/status'.format(domain),
                      lambda r: json.loads(r.text)['data']['status'] == 'Idle',
                      attempts=100)


def test_backup_rest(device, artifact_dir, domain):
    session = device.login_v2()
    response = session.post('https://{0}/rest/backup/create'.format(domain), json={'app': 'testapp'}, verify=False)
    assert response.status_code == 200
    assert json.loads(response.text)['success']

    wait_for_jobs(domain, session)

    response = device.http_get('/rest/backup/list')
    assert response.status_code == 200
    open('{0}/rest.backup.list.json'.format(artifact_dir), 'w').write(response.text)
    print(response.text)
    backup = json.loads(response.text)['data'][0]
    device.run_ssh('tar tvf {0}/{1}'.format(backup['path'], backup['file']))

    response = session.post(
        'https://{0}/rest/backup/restore'.format(domain),
        json={'app': 'testapp', 'file': '{0}'.format(backup['file'])},
        verify=False)
    assert response.status_code == 200

    wait_for_jobs(domain, session)


def test_backup_cli(device, artifact_dir):
    device.run_ssh("snap run platform.cli backup create testapp")
    response = device.run_ssh("snap run platform.cli backup list")
    open('{0}/cli.backup.list.json'.format(artifact_dir), 'w').write(response)
    print(response)
    backup = json.loads(response)[0]
    device.run_ssh('tar tvf {0}/{1}'.format(backup['path'], backup['file']))
    device.run_ssh("snap run platform.cli backup restore {0}".format(backup['file']))


def test_rest_backup_list(device, domain, artifact_dir):
    response = device.login_v2().get('https://{0}/rest/backup/list'.format(domain), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.backup.list.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)
    assert json.loads(response.text)['success']


@pytest.fixture(scope='function')
def loop_device(device_host):
    dev_file = '/tmp/disk'
    loop_device_cleanup(device_host, dev_file, password=LOGS_SSH_PASSWORD)

    print('adding loop device')
    run_ssh(device_host, 'dd if=/dev/zero bs=20M count=10 of={0}'.format(dev_file), password=LOGS_SSH_PASSWORD)
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


def disk_writable(domain):
    run_ssh(domain, 'ls -la /data/', password=LOGS_SSH_PASSWORD)
    run_ssh(domain, "touch /data/test.file", password=LOGS_SSH_PASSWORD)

mkfs = {
  'btrfs': '/snap/platform/current/btrfs/bin/mkfs.sh',
  'ext4': 'mkfs.ext4'
}
@pytest.mark.parametrize("fs_type", ['ext4', 'btrfs'])
def test_public_settings_disk_add_remove(loop_device, device, fs_type, domain, artifact_dir):
    disk_create(loop_device, fs_type, device)
    assert disk_activate(loop_device, device, domain, artifact_dir) == '/opt/disk/external/platform'
    disk_writable(domain)
    assert disk_deactivate(loop_device, device, domain) == '/opt/disk/internal/platform'


def test_public_settings_partition_add_remove(loop_device, device, domain, artifact_dir):
    device.run_ssh('/snap/platform/current/bin/disk_format.sh {0}'.format(loop_device), retries=3)
    partition = device.run_ssh('lsblk -pl -o NAME,TYPE {0} | grep part | head -1'.format(loop_device)).split()[0]
    session = device.login_v2()
    response = session.post('https://{0}/rest/storage/activate/partition'.format(domain), verify=False,
                            json={'device': partition, 'format': False})
    assert response.status_code == 200, response.text
    wait_for_jobs(domain, session)
    assert current_disk_link(device) == '/opt/disk/external/platform'
    disk_writable(domain)
    assert disk_deactivate(loop_device, device, domain) == '/opt/disk/internal/platform'


def disk_create(loop, fs, device):
    tmp_disk = '/tmp/test'
    device.run_ssh('{0} {1}'.format(mkfs[fs], loop), retries=3)

    device.run_ssh('rm -rf {0}'.format(tmp_disk))
    device.run_ssh('mkdir {0}'.format(tmp_disk))
    device.run_ssh('sync')

    device.run_ssh('mount {0} {1}'.format(loop, tmp_disk), retries=3)
    for mount in device.run_ssh('mount', debug=True).splitlines():
        if 'loop' in mount:
            print(mount)
    device.run_ssh('umount {0}'.format(loop))


def disk_activate(loop, device, domain, artifact_dir):
    session = device.login_v2()
    response = session.get('https://{0}/rest/storage/disks'.format(domain))
    print(response.text)
    with open('{0}/rest.storage.disks.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)

    assert loop in response.text
    assert response.status_code == 200

    response = session.post('https://{0}/rest/storage/activate/disk'.format(domain), verify=False,
                            json={'devices': [loop], 'format': True})
    assert response.status_code == 200, response.text
    wait_for_jobs(domain, session)

    return current_disk_link(device)


def disk_deactivate(loop, device, domain):
    response = device.login_v2().post('https://{0}/rest/storage/deactivate'.format(domain), verify=False,
                                   json={'device': loop})
    assert response.status_code == 200
    return current_disk_link(device)


def current_disk_link(device):
    return device.run_ssh("realpath /data/platform")


def test_if_cron_is_empty_after_install(device_host):
    cron_is_empty_after_install(device_host)


def cron_is_empty_after_install(device_host):
    crontab = run_ssh(device_host, "crontab -l", password=LOGS_SSH_PASSWORD, throw=False)
    assert 'no crontab for root' in crontab, crontab


def test_installer_version(domain, device, artifact_dir):
    response = device.login_v2().get('https://{0}/rest/installer/version'.format(domain), allow_redirects=False,
                                  verify=False)
    with open('{0}/rest.installer.version.json'.format(artifact_dir), 'w') as the_file:
        the_file.write(response.text)

    assert response.status_code == 200, response.text


def test_local_upgrade(app_archive_path, device_host):
    local_install(device_host, LOGS_SSH_PASSWORD, app_archive_path)


def test_reinstall_local_after_upgrade(app_archive_path, device_host):
    local_install(device_host, LOGS_SSH_PASSWORD, app_archive_path)


def test_nginx_performance(device_host):
    print(check_output('ab -c 1 -n 1000 https://{0}/ping'.format(device_host), shell=True).decode())


def test_cli_ipv4(device):
    ipv4 = device.run_ssh('snap run platform.cli ipv4'.format(TMP_DIR), throw=False)
    assert re.search(r"^\d*?\.\d*?\.\d*?\.\d*?$", ipv4) is not None


def test_nginx_plus_flask_performance(device_host):
    print(check_output('ab -c 1 -n 1000 https://{0}/rest/id'.format(device_host), shell=True).decode())


def test_cli_user_add_remove(device):
    device.run_ssh('snap run platform.cli user add testuser --password=testpassword123')
    device.run_ssh('snap run platform.cli user remove testuser')


def test_cli_cert_fake(device):
    cert_path = '/var/snap/platform/current/syncloud.crt'
    ca_path = '/var/snap/platform/current/syncloud.ca.crt'
    cert_serial_before = device.run_ssh('cat {} | openssl x509 -noout -serial'.format(cert_path)).strip()
    ca_serial_before = device.run_ssh('cat {} | openssl x509 -noout -serial'.format(ca_path)).strip()
    device.run_ssh('snap run platform.cli cert --fake')
    cert_serial_after = device.run_ssh('cat {} | openssl x509 -noout -serial'.format(cert_path)).strip()
    ca_serial_after = device.run_ssh('cat {} | openssl x509 -noout -serial'.format(ca_path)).strip()
    assert cert_serial_before != cert_serial_after, 'cert was not regenerated: {}'.format(cert_serial_after)
    assert ca_serial_before != ca_serial_after, 'CA was not regenerated: {}'.format(ca_serial_after)
    device.run_ssh('openssl verify -CAfile {} {}'.format(ca_path, cert_path))


def test_admin_api_secured(device, domain):
    # Unauthenticated requests to admin endpoints should return 401
    admin_endpoints = [
        '/rest/settings/2fa',
        '/rest/backup/list',
        '/rest/storage/disks',
        '/rest/certificate',
        '/rest/access',
        '/rest/apps/available',
    ]
    for endpoint in admin_endpoints:
        response = requests.get('https://{0}{1}'.format(domain, endpoint), verify=False, allow_redirects=False)
        assert response.status_code == 401, 'Expected 401 for {}, got {}'.format(endpoint, response.status_code)


def test_non_admin_api_allowed(device, domain):
    # Non-admin accessible endpoints should return 401 when unauthenticated (not 403)
    user_endpoints = [
        '/rest/apps/installed',
    ]
    for endpoint in user_endpoints:
        response = requests.get('https://{0}{1}'.format(domain, endpoint), verify=False, allow_redirects=False)
        assert response.status_code == 401, 'Expected 401 for {}, got {}'.format(endpoint, response.status_code)


def retry(method, retries=10):
    attempt = 0
    exception = None
    while attempt < retries:
        try:
            return method()
        except Exception as e:
            exception = e
            print('error (attempt {0}/{1}): {2}'.format(attempt + 1, retries, str(e)))
            time.sleep(5)
        attempt += 1
    raise exception

def wait_for_installer(web_session, host, attempts=60):
    is_running = True
    attempt = 0
    while is_running and attempt < attempts:
        try:
            response = web_session.get('https://{0}/rest/installer/status'.format(host), verify=False)
            if response.status_code == 200:
                status = json.loads(response.text)
                is_running = status['data']['is_running']
        except Exception as e:
            print(str(e))

        print("attempt: {0}/{1}".format(attempt, attempts))
        attempt += 1
        time.sleep(10)

    if is_running:
        raise Exception("time out waiting for thr installer")
