import sys
import logging

from syncloud_app import logger
from syncloud_platform.events import run_hook_script
from syncloud_platform.injector import get_injector

logger.init(logging.DEBUG, console=True, line_format='%(message)s')

apps_dir = sys.argv[1] # not used
app_id = sys.argv[2]
hook_script = sys.argv[3]

app_paths = get_injector().get_app_setup(app_id)
run_hook_script(app_paths, hook_script)
