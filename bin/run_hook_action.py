import sys
import logging

from syncloud_app import logger
from syncloud_platform.events import run_hook_script

logger.init(logging.DEBUG, console=True, line_format='%(message)s')

apps_dir = sys.argv[1]
app_id = sys.argv[2]
hook_script = sys.argv[3]

run_hook_script(apps_dir, app_id, hook_script)
