import os
import shutil
from os.path import dirname, join, exists

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


def test_web_with_selenium(user_domain):

    os.environ['PATH'] = os.environ['PATH'] + ":" + join(DIR, 'geckodriver')

    caps = DesiredCapabilities.FIREFOX
    caps["marionette"] = True
    caps["binary"] = "/usr/bin/firefox"
    caps['loggingPrefs'] = {'browser': 'ALL'}

    profile = webdriver.FirefoxProfile()
    profile.add_extension('JSErrorCollector.xpi')



    driver = webdriver.Firefox(profile, capabilities=caps, log_path="{0}/firefox.log".format(LOG_DIR))

    screenshot_dir = join(DIR, 'screenshot')
    if exists(screenshot_dir):
        shutil.rmtree(screenshot_dir)
    os.mkdir(screenshot_dir)

    driver.get("http://{0}:81".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(2)
    driver.get_screenshot_as_file(join(screenshot_dir, 'activate.png'))

    driver.get("http://{0}".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(2)
    driver.get_screenshot_as_file(join(screenshot_dir, 'login.png'))

    user = driver.find_element_by_id("name")
    user.send_keys(DEVICE_USER)
    password = driver.find_element_by_id("password")
    password.send_keys(DEVICE_PASSWORD)
    password.submit()
    wait_driver = WebDriverWait(driver, 10)
    wait_driver.until(EC.presence_of_element_located((By.CLASS_NAME, 'menubutton')))

    driver.get_screenshot_as_file(join(screenshot_dir, 'index.png'))

    driver.get("http://{0}/settings.html".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(2)
    driver.get_screenshot_as_file(join(screenshot_dir, 'settings.png'))

    driver.get("http://{0}/access.html".format(user_domain))
    wait_driver = WebDriverWait(driver, 10)
    time.sleep(5)
    driver.get_screenshot_as_file(join(screenshot_dir, 'access.png'))

    print(driver.execute_script('return window.JSErrorCollector_errors.pump()'))



