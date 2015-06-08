#!/usr/bin/env python
import logging
from os.path import dirname, join, abspath

from syncloud.app import logger
from syncloud.insider.facade import get_insider
from syncloud.sam.pip import Pip
from syncloud.sam.platform_installer import PlatformInstaller

APP_DIR = abspath(join(dirname(__file__), '..'))

logger.init(logging.DEBUG, True)

print("installing local binary build")
PlatformInstaller().install(join(APP_DIR, 'platform.tar.gz'))
Pip(None).log_version('syncloud-platform')

# persist upnp mock setting
get_insider().insider_config.set_upnpc_mock(True)
