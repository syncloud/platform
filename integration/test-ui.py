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


def test_internal_ui(driver, device_host, screenshot_dir):
    driver.get("http://{0}:81".format(device_host))
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'activate-old')

def test_activate(driver, device_host, screenshot_dir):
    driver.get("http://{0}/activate.html".format(device_host))
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'activate')

def test_login(driver, ui_mode, device_host, screenshot_dir):
    driver.get("http://{0}".format(device_host))
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'login-' + ui_mode)


def test_index(driver, ui_mode, device_user, device_password, screenshot_dir):
    user = driver.find_element_by_id("name")
    user.send_keys(device_user)
    password = driver.find_element_by_id("password")
    password.send_keys(device_password)
    password.submit()
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'index-progress-' + ui_mode)
    wait_driver = WebDriverWait(driver, 20)
    wait_driver.until(EC.presence_of_element_located((By.CLASS_NAME, 'menubutton')))
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'index-' + ui_mode)


def test_settings(driver, device_host, ui_mode, screenshot_dir):
    driver.get("http://{0}/settings.html".format(device_host))
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'settings-' + ui_mode)


def test_settings_activation(driver, device_host, ui_mode, screenshot_dir):
    driver.get("http://{0}/activation.html".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_activation-' + ui_mode)


def test_settings_network(driver, device_host, ui_mode, screenshot_dir):
    driver.get("http://{0}/network.html".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_network-' + ui_mode)

    driver.find_element_by_css_selector(".bootstrap-switch-id-tgl_external").click()
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'settings_network_external_access-' + ui_mode)


def test_settings_storage(driver, device_host, ui_mode, screenshot_dir):
    url = "http://{0}/storage.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_storage-' + ui_mode)


def test_settings_updates(driver, device_host, ui_mode, screenshot_dir):
    url = "http://{0}/updates.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_updates-' + ui_mode)


def test_settings_support(driver, device_host, ui_mode, screenshot_dir):
    url = "http://{0}/support.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_support-' + ui_mode)


def test_settings_backup(driver, ui_mode, device_host, screenshot_dir):
    url = "http://{0}/backup.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_backup-' + ui_mode)


def test_app_center(driver, ui_mode, device_host, screenshot_dir):
    url = "http://{0}/appcenter.html".format(device_host)
    resp = requests.get(url, verify=False)
    assert resp.status_code == 200
    driver.get(url)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'appcenter-' + ui_mode)


def test_installed_app(driver, device_host, ui_mode, screenshot_dir):
    driver.get("http://{0}/app.html?app_id=files".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'app_installed-' + ui_mode)


def test_not_installed_app(driver, device_host, ui_mode, screenshot_dir):
    driver.get("http://{0}/app.html?app_id=nextcloud".format(device_host))
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'app_not_installed-' + ui_mode)
