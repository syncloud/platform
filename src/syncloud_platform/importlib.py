from os.path import join, isdir
from os import listdir
import sys

def import_modules(lib_path):
    libs = [join(lib_path, item) for item in listdir(lib_path) if isdir(join(lib_path, item))]
    map(lambda x: sys.path.append(x), libs)


from os.path import dirname, join
this_path = dirname(__file__)
lib_path = join(this_path, '..', '..', 'lib')

import_modules(lib_path)
