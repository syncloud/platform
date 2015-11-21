# block to import lib folder
import sys
from os import listdir
from os.path import join, dirname, isdir, abspath

app_path = abspath(join(dirname(__file__), '..'))

lib_path = join(app_path, 'lib')
libs = [join(lib_path, item) for item in listdir(lib_path) if isdir(join(lib_path, item))]
map(lambda x: sys.path.append(x), libs)
# end of block to import lib folder


import syncloud_platform.rest.public
app = syncloud_platform.rest.public.app