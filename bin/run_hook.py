import logging

from syncloud_app import logger
from syncloud_platform.tools import app
from syncloud_platform.tools.nginx import Nginx
from syncloud_platform.systemd.systemctl import add_service, remove_service


import sys

script_filename = sys.argv[1]
logger.init(logging.DEBUG, console=True, line_format='%(message)s')
execfile(script_filename)