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

    os.environ['PATH'] = os.environ['PATH'] + ":" + join(DIR, 'geckodriver')

    caps = DesiredCapabilities.FIREFOX
    caps["marionette"] = True
    caps["binary"] = "/usr/bin/firefox"
    caps['loggingPrefs'] = {'browser': 'ALL'}

    profile = webdriver.FirefoxProfile()
    profile.add_extension('{0}/JSErrorCollector.xpi'.format(DIR))
    driver = webdriver.Firefox(profile, capabilities=caps, log_path="{0}/firefox.log".format(LOG_DIR))
    driver.set_page_load_timeout(30)
    return driver


@pytest.fixture(scope="module")
def module_setup(request):
    request.addfinalizer(module_teardown)


def module_teardown(driver):
    driver.close()
    

def test_internal_ui(driver, user_domain):

    driver.get("http://{0}:81".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'activate')
    print(driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []'))


def test_external_ui(driver, user_domain):

    driver.get("http://{0}".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(2)
    screenshots(driver, screenshot_dir, 'login')
    print(driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []'))


def test_login(driver, user_domain):

    user = driver.find_element_by_id("name")
    user.send_keys(DEVICE_USER)
    password = driver.find_element_by_id("password")
    password.send_keys(DEVICE_PASSWORD)
    password.submit()
    wait_driver = WebDriverWait(driver, 10)
    wait_driver.until(EC.presence_of_element_located((By.CLASS_NAME, 'menubutton')))
    screenshots(driver, screenshot_dir, 'index')

    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_settings(driver, user_domain):

    driver.get("http://{0}/settings.html".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(5)
    screenshots(driver, screenshot_dir, 'settings')
 
    assert not driver.execute_script('return window.JSErrorCollector_errors ? window.JSErrorCollector_errors.pump() : []')


def test_access(driver, user_domain):

    driver.get("http://{0}/access.html".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(10)
    screenshots(driver, screenshot_dir, 'access')
 
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

