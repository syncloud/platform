import logging

from syncloud.app import logger
from syncloud.tools import app
from syncloud.tools.nginx import Nginx
from syncloud.systemd.systemctl import add_service, remove_service


import sys

script_filename = sys.argv[1]
logger.init(logging.DEBUG, console=True, line_format='%(message)s')
execfile(script_filename)