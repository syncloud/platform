from os.path import dirname, join
from subprocess import check_output

import pytest
import time
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.screenshots import screenshots

DIR = dirname(__file__)
TMP_DIR = '/tmp/syncloud/ui'


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir, ui_mode, data_dir):
    def module_teardown():
        device.activated()
        ui_logs = join(artifact_dir, 'ui-log-{0}'.format(ui_mode))
        check_output('mkdir {0}'.format(ui_logs), shell=True)
        device.run_ssh('mkdir -p {0}'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/log/syslog {0}/syslog.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), ui_logs)
        device.scp_from_device('{0}/log/*'.format(data_dir), ui_logs)
        check_output('cp /videos/* {0}'.format(artifact_dir), shell=True)
        check_output('chmod -R a+r {0}'.format(ui_logs), shell=True)

    request.addfinalizer(module_teardown)


def test_start(app, device_host, module_setup, domain):
    add_host_alias(app, device_host, domain)


def test_deactivate(device, main_domain, domain):
    device.activated()
    device.run_ssh('snap run platform.cli config set redirect.domain {}'.format(main_domain))
    device.run_ssh('snap run platform.cli config set certbot.staging true')
    device.run_ssh('snap run platform.cli config set redirect.api_url http://api.redirect')

    response = device.login().post('https://{0}/rest/settings/deactivate'.format(domain), verify=False)
    assert '"success": true' in response.text
    assert response.status_code == 200


def test_fake_cert(selenium, device, device_host):
    device.run_ssh('rm /var/snap/platform/current/syncloud.crt')
    device.run_ssh('snap run platform.cli cert')
    device.run_ssh('snap restart platform')
    selenium.driver.get("https://{0}".format(device_host))
    selenium.screenshot('fake-cert')


def test_activate(selenium, device_host,
                  domain, device_user, device_password, redirect_user, redirect_password):
    selenium.driver.get("https://{0}".format(device_host))
    selenium.find_by_xpath("//h1[text()='Activate']")
    selenium.screenshot('activate-empty')
    selenium.find_by_id('btn_free_domain').click()
    wait_for(selenium, lambda: selenium.find_by_id('email').send_keys(""))
    selenium.find_by_id('email').send_keys(redirect_user)
    selenium.screenshot('activate-redirect-email')
    selenium.find_by_id('redirect_password').send_keys(redirect_password)
    selenium.find_by_id('domain_input').send_keys(domain)
    selenium.screenshot('activate-type')
    selenium.find_by_id('btn_next').click()
    selenium.screenshot('activate-redirect')
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.ID, 'device_username')))
    selenium.wait_or_screenshot(EC.presence_of_element_located((By.ID, 'device_password')))
    wait_for(selenium, lambda: selenium.find_by_id('device_username').send_keys(""))
    selenium.find_by_id('device_username').send_keys(device_user)
    selenium.find_by_id('device_password').send_keys(device_password)
    selenium.screenshot('activate-ready')
    selenium.find_by_id('btn_activate').click()
    selenium.find_by_xpath("//h1[text()='Log in']")


def test_activate_again(selenium, device_host):
    selenium.driver.get("https://{0}/activate".format(device_host))
    selenium.find_by_xpath("//h1[text()='Log in']")
    selenium.screenshot('activate')


def test_login(selenium, device_host):
    selenium.driver.get("https://{0}".format(device_host))
    selenium.find_by_xpath("//h1[text()='Log in']")
    selenium.screenshot('login')


def test_index(selenium, device_user, device_password):
    selenium.find_by_id("username").send_keys(device_user)
    selenium.find_by_id("password").send_keys(device_password)
    selenium.find_by_id("btn_login").click()
    selenium.screenshot('index-progress')
    selenium.find_by_xpath("//h1[text()='Applications']")
    selenium.screenshot('index')


def test_settings(selenium):
    menu(selenium, 'settings')
    selenium.find_by_xpath("//h1[text()='Settings']"((By.XPATH, header)))
    selenium.screenshot('settings')


def test_settings_activation(selenium):
    settings(selenium, 'activation')
    selenium.find_by_xpath("//h1[text()='Activation']")
    selenium.screenshot('settings_activation')


def test_settings_access(selenium):
    settings(selenium, 'access')
    selenium.find_by_xpath("//h1[text()='Access']")
    # selenium.find_by_xpath('//input[@id="tgl_ipv4_enabled"]/../span').click()
    selenium.find_by_xpath('//input[@id="tgl_ipv4_public"]/../span').click()
    # selenium.find_by_xpath('//input[@id="tgl_ip_autodetect"]/../span').click()
    # selenium.find_by_xpath('//input[@id="tgl_ipv6_enabled"]/../span').click()
    selenium.screenshot('settings_access')


def test_settings_network(selenium):
    settings(selenium, 'network')
    selenium.find_by_xpath("//h1[text()='Network']")
    selenium.screenshot('settings_network')


def test_settings_storage(selenium):
    settings(selenium, 'storage')
    selenium.find_by_id("//h1[text()='Storage']")
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
    selenium.find_by_id("auto").click()
    selenium.find_by_id("auto-backup").click()
    selenium.find_by_id("auto-day").click()
    selenium.find_by_id("auto-day-monday").click()
    selenium.find_by_id("auto-hour").click()
    selenium.find_by_id("auto-hour-1").click()
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


def test_installed_app(selenium):
    menu(selenium, 'appcenter')
    selenium.find_by_xpath("//h1[text()='App Center']")
    selenium.find_by_xpath("//span[text()='File browser']")
    selenium.find_by_xpath(files).click()
    selenium.find_by_xpath("//h1[text()='File browser']")
    selenium.screenshot('app_installed')


# def test_remove_app(driver, ui_mode, screenshot_dir):
#     remove = 'btn_remove'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, remove)))
#     driver.find_element_by_id(remove).click()
#     confirm = 'btn_confirm'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, confirm)))
#     driver.find_element_by_id(confirm).click()
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.invisibility_of_element_located((By.ID, remove)))
#     selenium.screenshot('app_removed-' + ui_mode)


# def test_install_app(driver, ui_mode, screenshot_dir):
#     install = 'btn_install'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.element_to_be_clickable((By.ID, install)))
#     driver.find_element_by_id(install).click()
#     confirm = 'btn_confirm'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, confirm)))
#     driver.find_element_by_id(confirm).click()
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.invisibility_of_element_located((By.ID, install)))
#     selenium.screenshot('app_installed-' + ui_mode)


def test_not_installed_app(selenium):
    menu(selenium, 'appcenter')
    selenium.find_by_xpath("//span[text()='Nextcloud file sharing']").click()
    selenium.find_by_xpath("//h1[text()='Nextcloud file sharing']")
    selenium.screenshot('app_not_installed')


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
            selenium.wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.element_to_be_clickable((By.ID, find_id)))
            selenium.selenium.screenshot(element_id + '-' + ui_mode)
            selenium.find_element_by_id(find_id).click()
            if selenium.ui_mode == "mobile":
                selenium.wait_or_screenshot(EC.invisibility_of_element_located((By.ID, find_id)))
            return
        except Exception as e:
            exception = e
            print('error (attempt {0}/{1}): {2}'.format(retry + 1, retries, str(e)))
            time.sleep(1)
        retry += 1
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
    selenium.find_by_id(setting).click()


def test_teardown(driver):
    driver.quit()
