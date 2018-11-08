import os
import shutil
import socket
import time
from os.path import join, dirname, isdir, split
from subprocess import check_output
from os import makedirs
import json

import jinja2
import pytest
import requests
from requests.adapters import HTTPAdapter
from syncloudlib.integration.installer import local_install, wait_for_sam, wait_for_rest, local_remove, \
    get_data_dir, get_app_dir, get_service_prefix, get_ssh_env_vars
from syncloudlib.integration.loop import loop_device_cleanup
from syncloudlib.integration.ssh import run_scp, run_ssh
from requests.packages.urllib3.exceptions import InsecureRequestWarning

requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

SYNCLOUD_INFO = 'syncloud.info'

DIR = dirname(__file__)
DEVICE_USER = "user"
DEVICE_PASSWORD = "password"
DEFAULT_DEVICE_PASSWORD = "syncloud"
LOGS_SSH_PASSWORD = DEFAULT_DEVICE_PASSWORD
LOG_DIR = join(DIR, 'log')


@pytest.fixture(scope="session")
def app_data_dir(installer):
    return get_data_dir(installer, 'app')

       
@pytest.fixture(scope="session")
def data_dir(installer):
    return get_data_dir(installer, 'platform')
         

@pytest.fixture(scope="session")
def app_dir(installer):
    return get_app_dir(installer, 'platform')


@pytest.fixture(scope="session")
def service_prefix(installer):
    return get_service_prefix(installer)


@pytest.fixture(scope="session")
def ssh_env_vars(installer):
    return get_ssh_env_vars(installer, 'platform')


@pytest.fixture(scope="session")
def module_setup(request, data_dir, device_host, app_dir):
    request.addfinalizer(lambda: module_teardown(data_dir, device_host, app_dir))


def module_teardown(data_dir, device_host, app_dir):
    run_scp('-r root@{0}:{1}/config {2}'.format(device_host, data_dir, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)
    run_scp('-r root@{0}:{1}/config.runtime {2}'.format(device_host, data_dir, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)

    run_ssh(device_host, 'journalctl > /log/journalctl.log', password=LOGS_SSH_PASSWORD, throw=False)
    run_scp('root@{0}:/log/* {1}'.format(device_host, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)
    run_scp('root@{0}:{1}/log/* {2}'.format(device_host, data_dir, LOG_DIR), throw=False, password=LOGS_SSH_PASSWORD)


def test_start(module_setup, device_host):
    shutil.rmtree(LOG_DIR, ignore_errors=True)
    run_scp('-r {0} root@{1}:/'.format(DIR, device_host))
    run_ssh(device_host, 'mkdir /log', password=LOGS_SSH_PASSWORD, throw=False)

    os.mkdir(LOG_DIR)


def test_install(app_archive_path, installer, device_host):
    local_install(device_host, DEFAULT_DEVICE_PASSWORD, app_archive_path, installer)


def test_http_port_validation_url(device_host):
    response = requests.get('http://{0}/ping'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert response.text == 'OK'


def test_https_port_validation_url(device_host):
    response = requests.get('https://{0}/ping'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 200
    assert response.text == 'OK'


def test_non_activated_device_main_page_redirect_to_activation(device_host):
    response = requests.get('https://{0}'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://{0}:81'.format(device_host)


def test_non_activated_device_login_redirect_to_activation(device_host):
    response = requests.post('https://{0}/rest/login'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 302
    assert response.headers['Location'] == 'http://{0}:81'.format(device_host)


def test_internal_web_open(device_host):

    response = requests.get('http://{0}:81'.format(device_host))
    assert response.status_code == 200


def test_activate_device(auth, device_host):

    email, password, domain, release = auth
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = 'password1'
    response = requests.post('http://{0}:81/rest/activate'.format(device_host),
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': 'user1', 'device_password': 'password1'})
    assert response.status_code == 200, response.text
    

def test_reactivate(auth, device_host):
    email, password, domain, release = auth
    response = requests.post('http://{0}:81/rest/activate'.format(device_host),
                             data={'main_domain': SYNCLOUD_INFO, 'redirect_email': email, 'redirect_password': password,
                                   'user_domain': domain, 'device_username': DEVICE_USER, 'device_password': DEVICE_PASSWORD})
    assert response.status_code == 200
    global LOGS_SSH_PASSWORD
    LOGS_SSH_PASSWORD = DEVICE_PASSWORD


def test_public_web_unauthorized_browser_redirect(device_host):
    response = requests.get('https://{0}/rest/user'.format(device_host), allow_redirects=False, verify=False)
    assert response.status_code == 302


def test_public_web_unauthorized_ajax_not_redirect(device_host):
    response = requests.get('https://{0}/rest/user'.format(device_host),
                            allow_redirects=False, verify=False, headers={'X-Requested-With': 'XMLHttpRequest'})
    assert response.status_code == 401


def test_running_platform_web(device_host):
    print(check_output('nc -zv -w 1 {0} 80'.format(device_host), shell=True))


def test_platform_rest(device_host):
    session = requests.session()
    session.mount('https://{0}'.format(device_host), HTTPAdapter(max_retries=5))
    response = session.get('https://{0}'.format(device_host), timeout=60, verify=False)
    assert response.status_code == 200


def test_app_unix_socket(app_dir, data_dir, app_data_dir, main_domain):
    nginx_template = '{0}/nginx.app.test.conf'.format(DIR)
    nginx_runtime = '{0}/nginx.app.test.conf.runtime'.format(DIR)
    generate_file_jinja(nginx_template, nginx_runtime, {'app_data': app_data_dir, 'platform_data': data_dir})
    run_scp('{0} root@{1}:/'.format(nginx_runtime, main_domain), throw=False, password=LOGS_SSH_PASSWORD)
    run_ssh(main_domain, 'mkdir -p {0}'.format(app_data_dir), password=DEVICE_PASSWORD)
    run_ssh(main_domain, '{0}/nginx/sbin/nginx '
                         '-c /nginx.app.test.conf.runtime '
                         '-g \'error_log {1}/log/test_nginx_app_error.log warn;\''.format(app_dir, data_dir),
            password=DEVICE_PASSWORD)
    response = requests.get('https://app.{0}'.format(main_domain), timeout=60, verify=False)
    assert response.status_code == 200
    assert response.text == 'OK', response.text


def test_api_service_restart(main_domain, ssh_env_vars):
    response = run_ssh(main_domain, '{0}/python/bin/python /integration/api_wrapper_service_restart.py platform.nginx-internal'.format(app_dir), password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
    assert 'OK' in response, response


def test_api_install_path(app_dir, main_domain, ssh_env_vars):
    response = run_ssh(main_domain, '{0}/python/bin/python /integration/api_wrapper_app_dir.py platform'.format(app_dir), password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
    assert app_dir in response, response
 
    
def test_api_data_path(app_dir, data_dir, main_domain, ssh_env_vars):
    response = run_ssh(main_domain, '{0}/python/bin/python /integration/api_wrapper_data_dir.py platform'.format(app_dir), password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
    assert data_dir in response, response
 
 
def test_api_url(app_dir, main_domain, user_domain, ssh_env_vars):
    response = run_ssh(main_domain, '{0}/python/bin/python /integration/api_wrapper_app_url.py platform'.format(app_dir), password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
    assert user_domain in response, response


def generate_file_jinja(from_path, to_path, variables):
    from_path_dir, from_path_filename = split(from_path)
    loader = jinja2.FileSystemLoader(searchpath=from_path_dir)

    env_parameters = dict(
        loader=loader,
        # some files like udev rules want empty lines at the end
        # trim_blocks=True,
        # lstrip_blocks=True,
        undefined=jinja2.StrictUndefined
    )
    environment = jinja2.Environment(**env_parameters)
    template = environment.get_template(from_path_filename)
    output = template.render(variables)
    to_path_dir = dirname(to_path)
    if not isdir(to_path_dir):
        makedirs(to_path_dir)
    with open(to_path, 'wb+') as fh:
        fh.write(output.encode("UTF-8"))


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


def test_set_access_mode_with_certbot(public_web_session, device_host):

    response = public_web_session.get('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                      params={'upnp_enabled': 'false',
                                              'external_access': 'false', 'public_ip': 0, 
                                              'certificate_port': 0, 'access_port': 0})
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_show_https_certificate(device_host):
    run_ssh(device_host, "echo | "
            "openssl s_client -showcerts -servername localhost -connect localhost:443 2>/dev/null | "
            "openssl x509 -inform pem -noout -text", password=DEVICE_PASSWORD)


def test_get_access(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/access/access'.format(device_host), verify=False)
    print(response.text)
    assert '"success": true' in response.text
    assert '"upnp_enabled": false' in response.text
    assert response.status_code == 200


def test_network_interfaces(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/access/network_interfaces'.format(device_host), verify=False,)
    print(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_available_apps(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/available_apps'.format(device_host), verify=False)
    with open('{0}/rest.available_apps.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)
    #assert '"success": true' in response.text
    assert response.status_code == 200
    assert len(json.loads(response.text)['apps']) > 1
    

def test_device_url(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/settings/device_url'.format(device_host), verify=False)
    with open('{0}/rest.settings.device_url.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_activate_url(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/settings/activate_url'.format(device_host), verify=False)
    with open('{0}/rest.settings.activate_url.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)
    assert '"success": true' in response.text
    assert response.status_code == 200


#def test_hook_override(public_web_session, data_dir, service_prefix, device_host):

    #run_ssh(device_host, "sed -i 's#hooks_root.*#hooks_root: /integration#g' {0}/config/platform.cfg".format(data_dir),
    #        password=DEVICE_PASSWORD)

    #run_ssh(device_host, 'systemctl restart {0}platform.uwsgi-public'.format(service_prefix), password=DEVICE_PASSWORD)

    #wait_for_rest(public_web_session, device_host, '/', 200)


def test_protocol(public_web_session, device_host, app_dir, ssh_env_vars, main_domain):

    response = public_web_session.get('https://{0}/rest/access/access'.format(device_host), verify=False)
    assert response.status_code == 200

    response = public_web_session.get('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                      params={'upnp_enabled': 'false',
                                              'external_access': 'false', 'public_ip': 0, 
                                              'certificate_port': 443, 'access_port': 443})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('https://{0}/rest/access/access'.format(device_host), verify=False)
    assert response.status_code == 200
    url = run_ssh(device_host, '{0}/python/bin/python /integration/api_wrapper_app_url.py platform'.format(app_dir), password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
   
    assert main_domain in url, url
    assert 'https' in url, url
   
    response = public_web_session.get('https://{0}/rest/access/set_access'.format(device_host), verify=False,
                                      params={'upnp_enabled': 'false',
                                              'external_access': 'false', 'public_ip': 0, 
                                              'certificate_port': 80, 'access_port': 10000})
    assert '"success": true' in response.text
    assert response.status_code == 200

    response = public_web_session.get('https://{0}/rest/access/access'.format(device_host), verify=False)
    assert response.status_code == 200

    url = run_ssh(device_host, '{0}/python/bin/python /integration/api_wrapper_app_url.py platform'.format(app_dir), password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
   
    assert main_domain in url, url
    assert 'https' in url, url
   

def test_sync(app_dir, ssh_env_vars, device_host):
    assert '"success": true' in run_ssh(device_host, '{0}/bin/insider sync_all'.format(app_dir),
                                        password=DEVICE_PASSWORD, env_vars=ssh_env_vars)
  
                                        
def test_cron(app_dir, ssh_env_vars, device_host):
    run_ssh(device_host, '{0}/bin/cron'.format(app_dir),
            password=DEVICE_PASSWORD, env_vars=ssh_env_vars)

                                        
def test_install_app(public_web_session, device_host):
    public_web_session.get('https://{0}/rest/install?app_id={1}'.format(device_host, 'files'), verify=False)
    wait_for_sam(public_web_session, device_host)


def test_rest_installed_apps(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/installed_apps'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.installed_apps.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)
    #assert '"success": true' in response.text
    assert response.status_code == 200
    assert len(json.loads(response.text)['apps']) == 1


def test_rest_installed_app(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/app?app_id=files'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.app.installed.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200
    #assert len(json.loads(response.text)['apps']) == 1

def test_rest_not_installed_app(public_web_session, device_host):
    response = public_web_session.get('https://{0}/rest/app?app_id=nextcloud'.format(device_host), verify=False)
    assert response.status_code == 200
    with open('{0}/rest.app.not.installed.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)
    assert response.status_code == 200
    #assert len(json.loads(response.text)['apps']) == 1


def test_do_not_cache_static_files_as_we_get_stale_ui_on_upgrades(public_web_session, device_host):

    response = public_web_session.get('https://{0}/settings.html'.format(device_host), verify=False)
    cache_control = response.headers['Cache-Control']
    assert 'no-cache' in cache_control
    assert 'max-age=0' in cache_control


def test_installer_upgrade(public_web_session, device_host):
    public_web_session.get('https://{0}/rest/upgrade?app_id=sam'.format(device_host), verify=False)
    wait_for_sam(public_web_session, device_host)


@pytest.yield_fixture(scope='function')
def loop_device(device_host, installer):
    dev_file = '/tmp/disk_{0}'.format(installer)
    loop_device_cleanup(device_host, dev_file, password=DEVICE_PASSWORD)

    print('adding loop device')
    run_ssh(device_host, 'dd if=/dev/zero bs=1M count=10 of={0}'.format(dev_file), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'sync', password=DEVICE_PASSWORD)
    run_ssh(device_host, 'ls -la {0}'.format(dev_file), password=DEVICE_PASSWORD)
    loop = run_ssh(device_host, 'losetup -f --show {0}'.format(dev_file), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'losetup', password=DEVICE_PASSWORD)
    run_ssh(device_host, 'losetup -j {0} | grep {0}'.format(dev_file), password=DEVICE_PASSWORD, retries=3)
    run_ssh(device_host, 'file -s {0}'.format(loop), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'sync', password=DEVICE_PASSWORD)
    run_ssh(device_host, 'partprobe {0}'.format(loop), password=DEVICE_PASSWORD, retries=3)
    yield loop

    loop_device_cleanup(device_host, dev_file, password=DEVICE_PASSWORD)


def disk_writable(device_host):
    run_ssh(device_host, 'ls -la /data/', password=DEVICE_PASSWORD)
    run_ssh(device_host, "touch /data/platform/test.file", password=DEVICE_PASSWORD)


@pytest.mark.parametrize("fs_type", ['ext4'])
def test_public_settings_disk_add_remove(loop_device, public_web_session, fs_type, device_host, installer, ssh_env_vars, app_dir):
    disk_create(loop_device, fs_type, device_host, installer)
    assert disk_activate(loop_device,  public_web_session, device_host, ssh_env_vars, app_dir) == '/opt/disk/external/platform'
    disk_writable(device_host)
    assert disk_deactivate(loop_device, public_web_session, device_host, ssh_env_vars, app_dir) == '/opt/disk/internal/platform'


def disk_create(loop_device, fs, device_host, installer):
    tmp_disk = '/tmp/test_{0}'.format(installer)
    run_ssh(device_host, 'mkfs.{0} {1}'.format(fs, loop_device), password=DEVICE_PASSWORD, retries=3)

    run_ssh(device_host, 'rm -rf {0}'.format(tmp_disk), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'mkdir {0}'.format(tmp_disk), password=DEVICE_PASSWORD)
    run_ssh(device_host, 'sync', password=DEVICE_PASSWORD)

    run_ssh(device_host, 'mount {0} {1}'.format(loop_device, tmp_disk), password=DEVICE_PASSWORD, retries=3)
    for mount in run_ssh(device_host, 'mount', debug=True, password=DEVICE_PASSWORD).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh(device_host, 'umount {0}'.format(loop_device), password=DEVICE_PASSWORD)


def disk_activate(loop_device, public_web_session, device_host, ssh_env_vars, app_dir):

    response = public_web_session.get('http://{0}/rest/settings/disks'.format(device_host))
    print response.text
    assert loop_device in response.text
    assert response.status_code == 200

    response = public_web_session.get('https://{0}/rest/settings/disk_activate'.format(device_host), verify=False,
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link(device_host, ssh_env_vars, app_dir)


def disk_deactivate(loop_device, public_web_session, device_host, ssh_env_vars, app_dir):
    response = public_web_session.get('https://{0}/rest/settings/disk_deactivate'.format(device_host), verify=False,
                                      params={'device': loop_device})
    assert response.status_code == 200
    return current_disk_link(device_host, ssh_env_vars, app_dir)


def current_disk_link(device_host, ssh_env_vars, app_dir):
    return run_ssh(device_host,
                   '{0}/python/bin/python /integration/api_wrapper_storage_init.py platform root'.format(app_dir),
                   password=DEVICE_PASSWORD, env_vars=ssh_env_vars)


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


def test_settings_versions(device_host, public_web_session):

    response = public_web_session.get('https://{0}/rest/settings/versions'.format(device_host), verify=False)
    with open('{0}/rest.settings.versions.json'.format(LOG_DIR), 'w') as the_file:
        the_file.write(response.text)

    assert response.status_code == 200, response.text


def test_local_upgrade(app_archive_path, installer, device_host):
    local_install(device_host, DEVICE_PASSWORD, app_archive_path, installer)


#def test_public_web_platform_upgrade(public_web_session, device_host, installer):
#    if installer == 'snapd':
#        run_ssh(device_host, 'snap refresh platform --amend --channel=stable', password=DEVICE_PASSWORD) 
        
#    public_web_session.get('https://{0}/rest/upgrade?app_id=platform&channel=master'.format(device_host), verify=False)
#    wait_for_sam(public_web_session, device_host)


def test_reinstall_local_after_upgrade(app_archive_path, installer, device_host):
    local_install(device_host, DEVICE_PASSWORD, app_archive_path, installer)


def test_if_cron_is_enabled_after_upgrade(device_host):
    cron_is_enabled_after_install(device_host)


def test_nginx_performance(device_host):
    print(check_output('ab -c 1 -n 1000 https://{0}/ping'.format(device_host), shell=True))


def test_nginx_plus_flask_performance(device_host):
    print(check_output('ab -c 1 -n 1000 http://{0}:81/rest/id'.format(device_host), shell=True))
