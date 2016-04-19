import sys
import logging
from syncloud_app import logger
from syncloud_platform.gaplib.scripts import run_script

logger.init(logging.DEBUG, console=True, line_format='%(message)s')
run_script(sys.argv[1])
