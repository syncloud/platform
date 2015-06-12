#!/usr/bin/env python
import logging
from os.path import dirname, join, abspath

from syncloud.app import logger
from syncloud.installer.installer import PlatformInstaller

APP_DIR = abspath(join(dirname(__file__), '..'))

logger.init(logging.DEBUG, True)

print("installing local binary build")
PlatformInstaller().install(join(APP_DIR, 'platform.tar.gz'))