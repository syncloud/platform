import os
import shutil
from os.path import dirname, join, exists
import pytest
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.firefox.firefox_binary import FirefoxBinary

DIR = dirname(__file__)
LOG_DIR = join(DIR, 'log')
DEVICE_USER = 'user'
DEVICE_PASSWORD = 'password'
screenshot_dir = join(DIR, 'screenshot')
 

@pytest.fixture(scope="module")
def driver():

    if exists(screenshot_dir):
        shutil.rmtree(screenshot_dir)
    os.mkdir(screenshot_dir)

    firefox_path = '/tools/firefox/firefox'
    caps = DesiredCapabilities.FIREFOX
    caps["marionette"] = True
    caps['acceptSslCerts'] = True

    binary = FirefoxBinary(firefox_path)

    profile = webdriver.FirefoxProfile()
    profile.add_extension('/tools/firefox/JSErrorCollector.xpi')
    profile.set_preference('app.update.auto', False)
    profile.set_preference('app.update.enabled', False)
    driver = webdriver.Firefox(profile,
                               capabilities=caps, log_path="{0}/firefox.log".format(LOG_DIR),
                               firefox_binary=binary, executable_path=join(DIR, '/tools/geckodriver/geckodriver'))

    #driver.set_page_load_timeout(30)
    #print driver.capabilities['version']
    return driver


@pytest.fixture(scope="module")
def module_setup(request):
    request.addfinalizer(module_teardown)


def module_teardown(driver):
    driver.close()
    

def test_internal_ui(driver, app_domain, device_host):

    driver.get("http://{0}:81".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'activate')
    print(driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []'))


def test_login(driver, app_domain, device_host):

    driver.get("http://{0}".format(device_host))
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'login')
    print(driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []'))

    
def test_index(driver, device_host):
    user = driver.find_element_by_id("name")
    user.send_keys(DEVICE_USER)
    password = driver.find_element_by_id("password")
    password.send_keys(DEVICE_PASSWORD)
    password.submit()
    wait_driver = WebDriverWait(driver, 10)
    wait_driver.until(EC.presence_of_element_located((By.CLASS_NAME, 'menubutton')))
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'index')

    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings(driver, device_host):
    
    driver.get("http://{0}/settings.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'settings')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_activation(driver, device_host):

    driver.get("http://{0}/activation.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_activation')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings_network(driver, device_host):

    driver.get("http://{0}/network.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_network')
    
    driver.find_element_by_id("tgl_external").click()
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'settings_network_external_access')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')

def test_settings_storage(driver, device_host):

    driver.get("http://{0}/storage.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_storage')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')

def test_settings_updates(driver, device_host):

    driver.get("http://{0}/updates.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_updates')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')

def test_settings_support(driver, device_host):

    driver.get("http://{0}/support.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'settings_support')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')

def test_app_center(driver, device_host):

    driver.get("http://{0}/appcenter.html".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'appcenter')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_installed_app(driver, device_host):

    driver.get("http://{0}/app.html?app_id=files".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'app_installed')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_not_installed_app(driver, device_host):

    driver.get("http://{0}/app.html?app_id=nextcloud".format(device_host))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'app_not_installed')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def screenshots(driver, dir, name):
    desktop_w = 1280
    desktop_h = 2000
    driver.set_window_position(0, 0)
    driver.set_window_size(desktop_w, desktop_h)

    driver.get_screenshot_as_file(join(dir, '{}.png'.format(name)))

    mobile_w = 400
    mobile_h = 2000
    driver.set_window_position(0, 0)
    driver.set_window_size(mobile_w, mobile_h)
    driver.get_screenshot_as_file(join(dir, '{}-mobile.png'.format(name)))
    
    driver.set_window_position(0, 0)
    driver.set_window_size(desktop_w, desktop_h)

