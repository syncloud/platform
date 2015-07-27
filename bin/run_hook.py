import logging
from syncloud_app import logger

# block to import lib folder
import sys
from os import listdir
from os.path import join, dirname, isdir, abspath

app_path = abspath(join(dirname(__file__), '..'))

lib_path = join(app_path, 'lib')
libs = [join(lib_path, item) for item in listdir(lib_path) if isdir(join(lib_path, item))]
map(lambda x: sys.path.insert(0, x), libs)
# end of block to import lib folder


script_filename = sys.argv[1]
logger.init(logging.DEBUG, console=True, line_format='%(message)s')
g = globals().copy()
g['__file__'] = script_filename
execfile(script_filename, g)
