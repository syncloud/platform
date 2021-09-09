import time
from subprocess import check_output
import pytest
from os.path import dirname, join
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from syncloudlib.integration.screenshots import screenshots
import requests

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
        check_output('chmod -R a+r {0}'.format(ui_logs), shell=True)

    request.addfinalizer(module_teardown)


def test_start(app, device_host, module_setup):
    pass


def test_deactivate(device, device_host):
    response = requests.get('https://{0}/rest/user'.format(device_host), allow_redirects=False, verify=False)
    if response.status_code != 501:
        response = device.login().post('https://{0}/rest/settings/deactivate'.format(device_host), verify=False)
        assert '"success": true' in response.text
        assert response.status_code == 200


def test_activate(driver, selenium, device_host,
                  domain, device_user, device_password, redirect_user, redirect_password):
    driver.get("https://{0}".format(device_host))
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
    #wait_for(selenium, lambda: selenium.find_by_id('device_username').send_keys(device_user))
    selenium.find_by_id('device_username').send_keys(device_user)
    selenium.find_by_id('device_password').send_keys(device_password)
    selenium.screenshot('activate-ready')
    selenium.find_by_id('btn_activate').click()
    selenium.find_by_xpath("//h1[text()='Log in']")


def test_activate_again(driver, ui_mode, device_host, screenshot_dir):
    driver.get("https://{0}/activate".format(device_host))
    header = "//h1[text()='Log in']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'activate')


def test_login(driver, ui_mode, device_host, screenshot_dir):
    driver.get("https://{0}".format(device_host))
    header = "//h1[text()='Log in']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'login-' + ui_mode)


def test_index(driver, ui_mode, device_user, device_password, screenshot_dir):
    user = driver.find_element_by_id("username")
    user.send_keys(device_user)
    password = driver.find_element_by_id("password")
    password.send_keys(device_password)
    login = driver.find_element_by_id("btn_login")
    login.click()
    screenshots(driver, screenshot_dir, 'index-progress-' + ui_mode)
    header = "//h1[text()='Applications']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'index-' + ui_mode)


def test_settings(driver, ui_mode, screenshot_dir):
    menu(driver, ui_mode, screenshot_dir, 'settings')
    header = "//h1[text()='Settings']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings-' + ui_mode)


def test_settings_activation(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'activation')
    header = "//h1[text()='Activation']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_activation-' + ui_mode)


def test_settings_access(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'access')
    header = "//h1[text()='Access']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    btn = 'external_mode'
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.element_to_be_clickable((By.ID, btn)))
    screenshots(driver, screenshot_dir, 'settings_access-' + ui_mode)
    driver.find_element_by_id(btn).click()
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, "tgl_ip_autodetect")))
    screenshots(driver, screenshot_dir, 'settings_access_external_access-' + ui_mode)


def test_settings_network(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'network')
    header = "//h1[text()='Network']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_network-' + ui_mode)


def test_settings_storage(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'storage')
    header = "//h1[text()='Storage']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_storage-' + ui_mode)


def test_settings_updates(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'updates')
    header = "//h1[text()='Updates']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_updates-' + ui_mode)


def test_settings_internal_memory(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'internalmemory')
    header = "//h1[text()='Internal Memory']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_updates-' + ui_mode)


def test_settings_support(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'support')
    header = "//h1[text()='Support']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_support-' + ui_mode)


def test_settings_backup(driver, ui_mode, screenshot_dir):
    settings(driver, screenshot_dir, ui_mode, 'backup')
    header = "//h1[text()='Backup']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'settings_backup-' + ui_mode)


def test_app_center(driver, ui_mode, screenshot_dir):
    menu(driver, ui_mode, screenshot_dir, 'appcenter')
    header = "//h1[text()='App Center']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    files = "//span[text()='File browser']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, files)))
    screenshots(driver, screenshot_dir, 'appcenter-' + ui_mode)


def test_installed_app(driver, ui_mode, screenshot_dir):
    menu(driver, ui_mode, screenshot_dir, 'appcenter')
    header = "//h1[text()='App Center']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    files = "//span[text()='File browser']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, files)))
    driver.find_element_by_xpath(files).click()
    header = "//h1[text()='File browser']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'app_installed-' + ui_mode)


# def test_remove_app(driver, ui_mode, screenshot_dir):
#     remove = 'btn_remove'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, remove)))
#     driver.find_element_by_id(remove).click()
#     confirm = 'btn_confirm'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, confirm)))
#     driver.find_element_by_id(confirm).click()
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.invisibility_of_element_located((By.ID, remove)))
#     screenshots(driver, screenshot_dir, 'app_removed-' + ui_mode)


# def test_install_app(driver, ui_mode, screenshot_dir):
#     install = 'btn_install'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.element_to_be_clickable((By.ID, install)))
#     driver.find_element_by_id(install).click()
#     confirm = 'btn_confirm'
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.ID, confirm)))
#     driver.find_element_by_id(confirm).click()
#     wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.invisibility_of_element_located((By.ID, install)))
#     screenshots(driver, screenshot_dir, 'app_installed-' + ui_mode)


def test_not_installed_app(driver, ui_mode, screenshot_dir):
    menu(driver, ui_mode, screenshot_dir, 'appcenter')
    nextcloud = "//span[text()='Nextcloud file sharing']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, nextcloud)))
    driver.find_element_by_xpath(nextcloud).click()
    header = "//h1[text()='Nextcloud file sharing']"
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.presence_of_element_located((By.XPATH, header)))
    screenshots(driver, screenshot_dir, 'app_not_installed-' + ui_mode)


def wait_or_screenshot(driver, ui_mode, screenshot_dir, method):
    wait_driver = WebDriverWait(driver, 120)
    try:
        wait_driver.until(method)
    except Exception as e:
        screenshots(driver, screenshot_dir, 'exception-' + ui_mode)
        raise e
    wait_for_loading(driver)


def menu(driver, ui_mode, screenshot_dir, element_id):
    wait_driver = WebDriverWait(driver, 30)
    retries = 10
    retry = 0
    exception = None
    while retry < retries:
        try:
            find_id = element_id
            if ui_mode == "mobile":
                find_id = element_id + '_mobile'
                menubutton = driver.find_element_by_id('menubutton')
                menubutton.click()
                wait_driver.until(EC.visibility_of_element_located((By.ID, find_id)))
            wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.element_to_be_clickable((By.ID, find_id)))
            screenshots(driver, screenshot_dir, element_id + '-' + ui_mode)
            element = driver.find_element_by_id(find_id)
            element.click()
            if ui_mode == "mobile":
                wait_driver.until(EC.invisibility_of_element_located((By.ID, find_id)))
            wait_for_loading(driver)
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


def settings(driver, screenshot_dir, ui_mode, setting):
    menu(driver, ui_mode, screenshot_dir, 'settings')
    wait_or_screenshot(driver, ui_mode, screenshot_dir, EC.element_to_be_clickable((By.ID, setting)))
    setting = driver.find_element_by_id(setting)
    setting.click()
    wait_for_loading(driver)


def wait_for_loading(driver):
    wait_driver = WebDriverWait(driver, 120)
    wait_driver.until(EC.invisibility_of_element_located((By.CLASS_NAME, 'loadingoverlay')))


def test_teardown(driver):
    driver.quit()

