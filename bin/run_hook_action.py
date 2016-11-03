import sys
import logging

from syncloud_app import logger
from syncloud_platform.events import run_hook_script
from syncloud_platform.config.config import PlatformConfig

logger.init(logging.DEBUG, console=True, line_format='%(message)s')

platform_config = PlatformConfig()
hook_script = sys.argv[1]
app_id = sys.argv[2]

run_hook_script(platform_config, hook_script, app_id)
