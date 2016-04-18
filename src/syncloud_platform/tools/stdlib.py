from os import makedirs
from os.path import isdir

def makepath(path):
    if not isdir(path):
        makedirs(path)