from os.path import dirname, join
from subprocess import check_output
from selenium.webdriver.support.ui import WebDriverWait

import pytest
import re
import socket
import time
import requests
import pyotp
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.http import wait_for_rest

DIR = dirname(__file__)
TMP_DIR = '/tmp/syncloud/ui'
stored_totp_secret = None


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir, ui_mode, data_dir):
    def module_teardown():
        device.activated()
        ui_logs = join(artifact_dir, ui_mode)
        check_output('mkdir -p {0}'.format(ui_logs), shell=True)
        device.run_ssh('mkdir -p {0}'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/platform/current/nginx.conf {0}/nginx.conf.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/platform/current/config/authelia/config.yml {0}/authelia.config.yml.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), ui_logs)
        check_output('cp /videos/* {0}'.format(artifact_dir), shell=True)
        check_output('chmod -R a+r {0}'.format(ui_logs), shell=True)

    request.addfinalizer(module_teardown)


def test_start(app, device_host, module_setup, domain, full_domain):
    add_host_alias(app, device_host, domain)
    add_host_alias("auth", device_host, full_domain)


def test_deactivate(device, device_host, main_domain, domain, full_domain):
    device.activated()
    ip = socket.gethostbyname(device_host)
    device.run_ssh('echo "{0} auth.{1}" >> /etc/hosts'.format(ip, full_domain))
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))
    device.run_ssh('snap run platform.cli config set certbot.staging true')
    device.run_ssh('snap run platform.cli config set redirect.api_url http://api.redirect')

    response = device.login_v2().post('https://{0}/rest/deactivate'.format(domain), verify=False)
    assert '"success":true' in response.text
    assert response.status_code == 200


def test_fake_cert(selenium, device, device_host):
    device.run_ssh('rm /var/snap/platform/current/syncloud.crt')
    device.run_ssh('snap run platform.cli cert')
    device.run_ssh('snap restart platform')
    wait_for_rest(requests.session(), "https://{0}/rest/activation/status".format(device_host), 200, 10)
    selenium.driver.get("https://{0}".format(device_host))
    selenium.find_by_xpath("//h1[text()='Activate']")
    wait_for_loading(selenium.driver)
    selenium.screenshot('fake-cert')


def test_activate(selenium, device_host,
                  domain, device_user, device_password, redirect_user, redirect_password):
    selenium.driver.get("https://{0}".format(device_host))
    selenium.find_by_xpath("//h1[text()='Activate']")
    selenium.screenshot('activate-empty')
    selenium.find_by_id('btn_free_domain').click()
    wait_for(selenium, lambda: selenium.find_by_id('email').send_keys(""))
    selenium.find_by_id('email').send_keys(redirect_user)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate-redirect-email')
    selenium.find_by_id('redirect_password').send_keys(redirect_password)
    selenium.find_by_id('domain_input').send_keys(domain)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate-type')
    selenium.find_by_id('btn_next').click()
    wait_for_loading(selenium.driver)
    selenium.screenshot('activate-redirect')
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.ID, 'device_username')))
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.ID, 'device_password')))
    wait_for(selenium, lambda: selenium.find_by_id('device_username').send_keys(""))
    selenium.find_by_id('device_username').send_keys(device_user)
    selenium.find_by_id('device_password').send_keys(device_password)
    selenium.find_by_id('device_password_confirm').send_keys(device_password)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate-ready')
    selenium.find_by_id('btn_activate').click()
    wait_for_loading(selenium.driver)
    selenium.find_by(By.ID, "username-textfield")
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate')


def test_login(selenium, full_domain, device_user, device_password):
    selenium.driver.get("https://{0}".format(full_domain))
    # OIDC flow redirects to Authelia
    selenium.find_by(By.ID, "username-textfield").send_keys(device_user)
    selenium.find_by(By.ID, "password-textfield").send_keys(device_password)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('login')
    selenium.find_by(By.ID, "sign-in-button").click()
    selenium.find_by_xpath("//h1[text()='Applications']")
    wait_for_loading(selenium.driver)
    selenium.find_by(By.CLASS_NAME, "appimg")
    selenium.screenshot('index')


def test_settings(selenium):
    menu(selenium, 'settings')
    selenium.find_by_xpath("//h1[text()='Settings']")
    selenium.screenshot('settings')


def test_settings_activation(selenium):
    settings(selenium, 'activation')
    selenium.find_by_xpath("//h1[text()='Activation']")
    selenium.screenshot('settings_activation')


def test_settings_access(selenium):
    settings(selenium, 'access')
    selenium.find_by_xpath("//h1[text()='Access']")
    selenium.find_by_xpath('//input[@id="tgl_ipv4_enabled"]/../span').click()
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.CSS_SELECTOR, '#ipv4_mode_block[data-ready]')))
    selenium.find_by_xpath('//input[@id="tgl_ipv4_public"]/../span').click()
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.CSS_SELECTOR, '#ipv4_public_block[data-ready]')))
    selenium.screenshot('settings_access')


def test_settings_network(selenium):
    settings(selenium, 'network')
    selenium.find_by_xpath("//h1[text()='Network']")
    selenium.screenshot('settings_network')


def test_settings_storage(selenium):
    settings(selenium, 'storage')
    selenium.find_by_xpath("//h1[text()='Storage']")
    selenium.find_by_id('btn_save')
    selenium.screenshot('settings_storage')


def test_settings_updates(selenium):
    settings(selenium, 'updates')
    selenium.find_by_xpath("//h1[text()='Updates']")
    selenium.screenshot('settings_updates')


def test_settings_internal_memory(selenium):
    settings(selenium, 'internalmemory')
    selenium.find_by_xpath("//h1[text()='Internal Memory']")
    selenium.screenshot('settings_internal_memory')


def test_settings_support(selenium):
    settings(selenium, 'support')
    selenium.find_by_xpath("//h1[text()='Support']")
    selenium.screenshot('settings_support')


def test_settings_backup(selenium):
    settings(selenium, 'backup')
    selenium.find_by_xpath("//h1[text()='Backup']")
    selenium.screenshot('settings_backup')
    assert not selenium.exists_by(By.CSS_SELECTOR, '.el-notification__title')
    selenium.clickable_by(By.ID, "auto").click()
    selenium.clickable_by(By.ID, "auto-backup").click()
    selenium.clickable_by(By.ID, "auto-day").click()
    selenium.clickable_by(By.ID, "auto-day-monday").click()
    selenium.clickable_by(By.ID, "auto-hour").click()
    selenium.clickable_by(By.ID, "auto-hour-1").click()
    selenium.find_by_id("save").click()
    selenium.screenshot('settings_backup_saved')
    assert not selenium.exists_by(By.CSS_SELECTOR, '.el-notification__title')


def test_settings_certificate(selenium):
    settings(selenium, 'certificate')
    selenium.find_by_xpath("//h1[text()='Certificate']")
    selenium.screenshot('settings_certificate')


def test_app_center(selenium):
    menu(selenium, 'appcenter')
    selenium.find_by_xpath("//h1[text()='App Center']")
    selenium.find_by_xpath("//span[text()='File browser']")
    selenium.screenshot('appcenter')


def test_custom_proxy_overrides_missing_app(selenium, device, device_host, full_domain):
    device.run_ssh('nohup /test/externalapp/externalapp > /tmp/syncloud/ui/externalapp.log 2>&1 &', throw=False)
    add_host_alias("files", device_host, full_domain)
    settings(selenium, 'customproxy')
    selenium.find_by_xpath("//h1[text()='Custom Proxy']")
    selenium.find_by_id('proxy_name').send_keys('files')
    selenium.find_by_id('proxy_host').send_keys('localhost')
    selenium.find_by_id('proxy_port').send_keys('8585')
    selenium.find_by_id('btn_add').click()
    wait_for_loading(selenium.driver)
    selenium.screenshot('settings_custom_proxy_files_added')

    def check_proxy():
        response = requests.get('https://files.{0}'.format(full_domain), verify=False)
        assert response.status_code == 200, response.text
        assert response.text == "external", response.text
    wait_for(selenium, check_proxy)
    selenium.screenshot('settings_custom_proxy_files_verified')


def test_install_app(selenium, device, full_domain):
    menu(selenium, 'appcenter')
    selenium.find_by_xpath("//h1[text()='App Center']")
    selenium.find_by_xpath("//span[text()='File browser']").click()
    selenium.find_by_xpath("//h1[text()='File browser']")
    selenium.find_by_id('btn_install').click()
    selenium.find_by_id('btn_confirm').click()
    # Store install takes several minutes, poll for btn_remove
    for attempt in range(60):
        if selenium.exists_by(By.ID, 'btn_remove'):
            break
        print('waiting for install to finish (attempt {0}/60)'.format(attempt + 1))
    selenium.find_by_id('btn_remove')
    selenium.screenshot('app_installed')

    # Debug: check what socket files the installed app created
    device.run_ssh('ls -la /var/snap/files/common/', throw=False)

    # After installing the real files app, the socket file exists
    # so nginx should route to the real app, not the custom proxy
    def check_real_app():
        response = requests.get('https://files.{0}'.format(full_domain), verify=False)
        assert response.status_code == 200, response.text
        assert response.text != "external", "custom proxy should not be used when app is installed"
    wait_for(selenium, check_real_app)


def test_remove_app(selenium):
    selenium.find_by_id('btn_remove').click()
    selenium.find_by_id('btn_confirm').click()
    wait_for_loading(selenium.driver)
    selenium.find_by_id("btn_install")
    selenium.screenshot('app_removed')


def test_remove_custom_proxy_files(selenium):
    settings(selenium, 'customproxy')
    selenium.find_by_xpath("//h1[text()='Custom Proxy']")
    selenium.find_by_id('btn_remove_files').click()
    wait_for_loading(selenium.driver)
    selenium.screenshot('settings_custom_proxy_files_removed')


def test_not_installed_app(selenium):
    menu(selenium, 'appcenter')
    selenium.clickable_by(By.XPATH, "//span[text()='Nextcloud file sharing']").click()
    selenium.find_by_xpath("//h1[text()='Nextcloud file sharing']")
    selenium.wait_or_screenshot(EC.visibility_of_element_located((By.ID, 'btn_install')))
    selenium.screenshot('app_not_installed')


def test_settings_custom_proxy(selenium, device, device_host, full_domain):
    add_host_alias("externalapp", device_host, full_domain)
    settings(selenium, 'customproxy')
    selenium.find_by_xpath("//h1[text()='Custom Proxy']")
    wait_for_loading(selenium.driver)
    selenium.screenshot('settings_custom_proxy')
    selenium.find_by_id('proxy_name').send_keys('externalapp')
    selenium.find_by_id('proxy_host').send_keys('localhost')
    selenium.find_by_id('proxy_port').send_keys('8585')
    selenium.screenshot('settings_custom_proxy_filled')
    selenium.find_by_id('btn_add').click()
    wait_for_loading(selenium.driver)
    selenium.find_by_xpath("//a[text()='externalapp']")
    selenium.screenshot('settings_custom_proxy_added')

    def check_proxy():
        response = requests.get('https://externalapp.{0}'.format(full_domain), verify=False)
        assert response.status_code == 200, response.text
        assert response.text == "external", response.text
    wait_for(selenium, check_proxy)

    selenium.screenshot('settings_custom_proxy_verified')


def logout(selenium, full_domain):
    selenium.driver.get("https://{0}/rest/logout".format(full_domain))
    # Wait for Authelia logout to complete (redirects to login form)
    selenium.find_by(By.ID, "username-textfield")


def test_auth_web(selenium, full_domain, device_user, device_password):
    logout(selenium, full_domain)
    selenium.driver.get("https://auth.{0}".format(full_domain))
    selenium.find_by(By.ID, "username-textfield").send_keys(device_user)
    password = selenium.find_by(By.ID, "password-textfield")
    password.send_keys(device_password)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('auth')
    selenium.find_by(By.ID, "sign-in-button").click()

    # redirect to the main web
    selenium.find_by_xpath("//h1[text()='Applications']")
    wait_for_loading(selenium.driver)


def test_2fa_settings(selenium, full_domain):
    selenium.driver.get("https://{0}".format(full_domain))
    settings(selenium, 'twofactor')
    selenium.find_by_xpath("//h1[text()='Two-Factor Authentication']")
    selenium.find_by_id('twofa_status')
    selenium.screenshot('2fa_settings')


def test_2fa_enable(selenium, device, full_domain, device_user, device_password):
    selenium.driver.get("https://{0}".format(full_domain))
    settings(selenium, 'twofactor')
    selenium.find_by_xpath("//h1[text()='Two-Factor Authentication']")

    selenium.find_by_id('btn_enable_2fa').click()
    # Wait for QR code to appear (backend enables 2FA + generates TOTP via CLI)
    for attempt in range(30):
        if selenium.exists_by(By.ID, 'totp_qr'):
            break
        print('waiting for TOTP QR code (attempt {0}/30)'.format(attempt + 1))
        time.sleep(2)
    selenium.find_by_id('totp_qr')
    selenium.screenshot('2fa_enabled')

    global stored_totp_secret
    secret_element = selenium.find_by_id('totp_secret')
    stored_totp_secret = secret_element.text
    selenium.find_by_id('btn_disable_2fa')
    selenium.screenshot('2fa_totp_registered')


def test_2fa_login(selenium, device, full_domain, device_user, device_password):
    logout(selenium, full_domain)
    selenium.driver.get("https://{0}".format(full_domain))

    selenium.find_by(By.ID, "username-textfield").send_keys(device_user)
    selenium.find_by(By.ID, "password-textfield").send_keys(device_password)
    selenium.find_by(By.ID, "sign-in-button").click()

    # TOTP challenge
    selenium.find_by(By.ID, "otp-input")
    selenium.screenshot('2fa_login_totp')
    totp = pyotp.TOTP(stored_totp_secret)
    # Wait for next TOTP period to avoid replay rejection
    remaining = totp.interval - time.time() % totp.interval
    time.sleep(remaining + 1)
    code = totp.now()
    otp_inputs = selenium.driver.find_elements(By.CSS_SELECTOR, "#otp-input input")
    for i, digit in enumerate(code):
        otp_inputs[i].send_keys(digit)

    selenium.find_by_xpath("//h1[text()='Applications']")
    wait_for_loading(selenium.driver)
    selenium.screenshot('2fa_login_success')


def test_2fa_disable(selenium, full_domain):
    settings(selenium, 'twofactor')
    selenium.find_by_xpath("//h1[text()='Two-Factor Authentication']")
    selenium.find_by_id('btn_disable_2fa').click()
    selenium.screenshot('2fa_disabled')


def test_2fa_recovery_cli(device, selenium, full_domain, device_user, device_password):
    # enable 2FA again via API
    session = device.login_v2()
    session.post('https://{0}/rest/settings/2fa'.format(full_domain),
                 json={'enabled': True}, verify=False)

    # disable via CLI
    device.run_ssh('snap run platform.cli disable-2fa')
    time.sleep(2)

    # verify login works without TOTP
    logout(selenium, full_domain)
    selenium.driver.get("https://{0}".format(full_domain))
    selenium.find_by(By.ID, "username-textfield").send_keys(device_user)
    selenium.find_by(By.ID, "password-textfield").send_keys(device_password)
    selenium.find_by(By.ID, "sign-in-button").click()
    selenium.find_by_xpath("//h1[text()='Applications']")
    wait_for_loading(selenium.driver)
    selenium.screenshot('2fa_recovery_cli')
 

def test_settings_deactivate(selenium, device_host, full_domain,
                  domain, device_user, device_password, redirect_user, redirect_password):
    selenium.driver.get("https://{0}".format(full_domain))
    settings(selenium, 'activation')
    selenium.find_by_xpath("//h1[text()='Activation']")
    selenium.find_by_id('btn_reactivate').click()
    selenium.find_by_xpath("//h1[text()='Activate']")
    selenium.screenshot('activate-empty')
    selenium.find_by_id('btn_free_domain').click()
    wait_for(selenium, lambda: selenium.find_by_id('email').send_keys(""))
    selenium.find_by_id('email').send_keys(redirect_user)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate-redirect-email')
    selenium.find_by_id('redirect_password').send_keys(redirect_password)
    selenium.find_by_id('domain_input').send_keys(domain)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate-type')
    selenium.find_by_id('btn_next').click()
    wait_for_loading(selenium.driver)
    selenium.screenshot('activate-redirect')
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.ID, 'device_username')))
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.ID, 'device_password')))
    wait_for(selenium, lambda: selenium.find_by_id('device_username').send_keys(""))
    selenium.find_by_id('device_username').send_keys(device_user)
    selenium.find_by_id('device_password').send_keys(device_password)
    selenium.find_by_id('device_password_confirm').send_keys(device_password)
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('activate-ready')
    selenium.find_by_id('btn_activate').click()
    wait_for_loading(selenium.driver)
    selenium.find_by(By.ID, "username-textfield")
    selenium.driver.find_element(By.TAG_NAME, "body").click()
    time.sleep(1)
    selenium.screenshot('deactivate-login-page')
    selenium.find_by(By.ID, "username-textfield").send_keys(device_user)
    selenium.find_by(By.ID, "password-textfield").send_keys(device_password)
    selenium.find_by(By.ID, "sign-in-button").click()
    selenium.find_by_xpath("//h1[text()='Applications']")
    wait_for_loading(selenium.driver)
    selenium.screenshot('reactivate-index')


def test_permission_denied(selenium, device, ui_mode, full_domain):
    device.run_ssh('/snap/platform/current/openldap/bin/ldapadd.sh -x -w syncloud -D "dc=syncloud,dc=org" -f /test/test.{0}.ldif'.format(ui_mode))
    logout(selenium, full_domain)
    selenium.driver.get("https://{0}".format(full_domain))
    selenium.find_by(By.ID, "username-textfield").send_keys("test{0}".format(ui_mode))
    selenium.find_by(By.ID, "password-textfield").send_keys("password")
    selenium.find_by(By.ID, "sign-in-button").click()
    selenium.find_by(By.CSS_SELECTOR, ".notification")
    selenium.wait_or_screenshot(EC.element_to_be_clickable((By.ID, "sign-in-button")))
    selenium.screenshot('permission-denied')


def test_502(selenium, full_domain):
    selenium.driver.get("https://unknown.{0}".format(full_domain))
    selenium.find_by_xpath("//h2[contains(.,'App is not available')]")


def menu(selenium, element_id):
    retries = 10
    retry = 0
    exception = None
    while retry < retries:
        try:
            find_id = element_id
            if selenium.ui_mode == "mobile":
                find_id = element_id + '_mobile'
                selenium.find_by_id('menubutton').click()
                selenium.wait_or_screenshot(EC.visibility_of_element_located((By.ID, find_id)))
            selenium.wait_or_screenshot(EC.element_to_be_clickable((By.ID, find_id)))
            selenium.find_by_id(find_id).click()
            if selenium.ui_mode == "mobile":
                wait_for_menu_close(selenium.driver)
            wait_for_loading(selenium.driver)
            selenium.screenshot(element_id)
            return
        except Exception as e:
            exception = e
            print('error (attempt {0}/{1}): {2}'.format(retry + 1, retries, str(e)))
            time.sleep(1)
        retry += 1
    selenium.screenshot('exception')
    raise exception


def wait_for(selenium, method):
    retries = 10
    retry = 0
    exception = None
    while retry < retries:
        try:
            method()
            return
        except Exception as e:
            exception = e
            print('error (attempt {0}/{1}): {2}'.format(retry + 1, retries, str(e)))
            time.sleep(1)
        retry += 1
    selenium.screenshot('exception')
    raise exception


def settings(selenium, setting):
    menu(selenium, 'settings')
    selenium.clickable_by(By.ID, setting).click()
    wait_for_loading(selenium.driver)


def wait_for_menu_close(driver):
    wait_driver = WebDriverWait(driver, 10)
    wait_driver.until(EC.invisibility_of_element_located((By.ID, 'menu')))


def wait_for_loading(driver):
    # Brief pause to let async loading overlays appear before checking invisibility
    time.sleep(0.5)
    wait_driver = WebDriverWait(driver, 120)
    wait_driver.until(EC.invisibility_of_element_located((By.CLASS_NAME, 'el-loading-mask')))
