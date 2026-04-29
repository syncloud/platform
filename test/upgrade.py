import json
from subprocess import run

import pytest
import requests
from syncloudlib.http import wait_for_rest, wait_for_response
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.installer import local_install

TMP_DIR = '/tmp/syncloud'
DB_PATH = '/var/snap/platform/current/platform.db'


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir):
    def module_teardown():
        device.run_ssh('journalctl > {0}/upgrade.journalctl.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), artifact_dir)
        run('cp /videos/* {0}'.format(artifact_dir), shell=True)
        run('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, app, device_host, domain, device):
    add_host_alias(app, device_host, domain)
    device.activated()
    device.run_ssh('rm -rf {0}'.format(TMP_DIR), throw=False)
    device.run_ssh('mkdir {0}'.format(TMP_DIR), throw=False)


def test_upgrade(device, device_user, device_password, device_host, app_archive_path, app_domain, app_dir):
    device.run_ssh('snap remove platform')
    device.run_ssh('/test/install-snapd.sh')
    device.run_ssh('snap install platform', retries=3)

    # Insert a custom proxy entry using old schema (without https column) to test migration
    # TODO: replace sqlite3 manipulation with CLI after this version goes to stable
    run('rm -rf /tmp/upgrade && mkdir /tmp/upgrade', shell=True, check=True)
    device.scp_from_device(DB_PATH, '/tmp/upgrade')
    run('sqlite3 /tmp/upgrade/platform.db "INSERT INTO custom_proxy (name, host, port) VALUES (\'testproxy\', \'localhost\', 8080)"',
        shell=True, check=True)
    device.scp_to_device('/tmp/upgrade/platform.db', DB_PATH)

    local_install(device_host, device_password, app_archive_path)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 10)


def test_activate_after_upgrade(device, device_host, device_user, device_password):
    wait_for_rest(requests.session(), "https://{0}/rest/activation/status".format(device_host), 200, 10)
    response = requests.post('https://{0}/rest/activate/custom'.format(device_host),
                             json={'domain': 'example.com',
                                   'device_username': device_user,
                                   'device_password': device_password}, verify=False)
    assert response.status_code == 200, response.text


def test_custom_proxy_migration(device):
    # TODO: replace sqlite3 manipulation with CLI after this version goes to stable
    output = device.run_ssh('snap run platform.cli proxy list')
    proxies = json.loads(output)
    assert len(proxies) == 1
    assert proxies[0]["name"] == "testproxy"
    assert proxies[0]["host"] == "localhost"
    assert proxies[0]["port"] == 8080
    assert proxies[0]["https"] is False


def test_slapd_quiet_after_refresh(device):
    output = device.run_ssh(
        "journalctl -u snap.platform.openldap --since '5 minutes ago' "
        "| grep -c 'slapd.* BIND ' || true"
    )
    assert int(output.strip()) == 0, output


def test_installer_upgrade(device, domain):
    session = device.login_v2()
    response = session.post('https://{0}/rest/installer/upgrade'.format(domain), verify=False)
    assert response.status_code == 200, response.text
    wait_for_jobs(domain, session)

    response = session.post('https://{0}/rest/installer/upgrade'.format(domain), verify=False)
    assert response.status_code == 200, response.text
    wait_for_jobs(domain, session)


def wait_for_jobs(domain, session):
    wait_for_response(session, 'https://{0}/rest/job/status'.format(domain),
                      lambda r: json.loads(r.text)['data']['status'] == 'Idle',
                      attempts=100)
