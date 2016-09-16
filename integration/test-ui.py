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
log_dir = join(LOG_DIR, 'log')


def test_web_with_selenium(user_domain):

    os.environ['PATH'] = os.environ['PATH'] + ":" + join(DIR, 'geckodriver')

    caps = DesiredCapabilities.FIREFOX
    caps["marionette"] = True
    caps["binary"] = "/usr/bin/firefox"

    profile = webdriver.FirefoxProfile()
    profile.set_preference("webdriver.log.file", "{0}/firefox.log".format(log_dir))
    driver = webdriver.Firefox(profile, capabilities=caps)

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



