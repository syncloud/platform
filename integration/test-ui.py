import os
import shutil
import time
from os.path import dirname, join, exists

import pytest
import requests
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from syncloudlib.integration.screenshots import screenshots

DIR = dirname(__file__)
LOG_DIR = join(DIR, 'log')
DEVICE_USER = 'user'
DEVICE_PASSWORD = 'password'
screenshot_dir = join(DIR, 'screenshot')


@pytest.fixture(scope="module")
def module_setup(request):
    request.addfinalizer(module_teardown)


def test_start():
    if exists(screenshot_dir):
        shutil.rmtree(screenshot_dir)
    os.mkdir(screenshot_dir)


def module_teardown(driver, mobile_driver):
    driver.close()
    mobile_driver.close()


def test_internal_ui(driver, device_host):
    driver.get("http://{0}:81".format(device_host))
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'activate')
    print(driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []'))


def teat_lpgin(driver, mobile_driver, device_host):
    _test_login(driver, 'desktop', device_host)
    _test_login(mobile_driver, 'mobile', device_host)


def _test_login(driver, mode, device_host):
    driver.get("http://{0}".format(device_host))
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'login-' + mode)
    print(driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []'))


def test_index(driver, mobile_driver):
    _test_index(driver, 'desktop')
    _test_index(mobile_driver, 'mobile')


def _test_index(driver, mode):
    user = driver.find_element_by_id("name")
    user.send_keys(DEVICE_USER)
    password = driver.find_element_by_id("password")
    password.send_keys(DEVICE_PASSWORD)
    password.submit()
    wait_driver = WebDriverWait(driver, 10)
    wait_driver.until(EC.presence_of_element_located((By.CLASS_NAME, 'menubutton')))
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'index-' + mode)

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings(driver, device_host):
    driver.get("http://{0}/settings.html".format(device_host))
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'settings')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_activation(driver, device_host):
    driver.get("http://{0}/activation.html".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_activation')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_network(driver, device_host):
    driver.get("http://{0}/network.html".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_network')

    driver.find_element_by_css_selector(".bootstrap-switch-id-tgl_external").click()
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'settings_network_external_access')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_storage(driver, device_host):
    url = "http://{0}/storage.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_storage')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_updates(driver, device_host):
    url = "http://{0}/updates.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_updates')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_support(driver, device_host):
    url = "http://{0}/support.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_support')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_app_center(driver, mobile_driver, device_host):
    _test_app_center(driver, 'desktop', device_host)
    _test_app_center(mobile_driver, 'mobile', device_host)


def _test_app_center(driver, mode, device_host):
    url = "http://{0}/appcenter.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'appcenter-' + mode)

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')
    assert not driver.execute_script(
        'return document.documentElement.scrollWidth>document.documentElement.clientWidth;')


def test_installed_app(driver, device_host):
    driver.get("http://{0}/app.html?app_id=files".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'app_installed')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_not_installed_app(driver, device_host):
    driver.get("http://{0}/app.html?app_id=nextcloud".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'app_not_installed')

    assert not driver.execute_script(
        'return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')
