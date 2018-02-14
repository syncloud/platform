import sys
import logging
from os.path import join
from syncloud_app import logger
from syncloud_platform.events import run_hook_script

logger.init(logging.DEBUG, console=True, line_format='%(message)s')

apps_dir = sys.argv[1]
app_id = sys.argv[2]
hook_script = sys.argv[3]
app_dir = join(apps_dir, app_id)
run_hook_script(app_dir, hook_script)
