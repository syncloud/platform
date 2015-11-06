import logging
from syncloud_app import logger


def run_script(script_filename):
    logger.init(logging.DEBUG, console=True, line_format='%(message)s')
    g = globals().copy()
    g['__file__'] = script_filename
    execfile(script_filename, g)
