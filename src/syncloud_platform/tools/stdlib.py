from os import makedirs
from os.path import isdir
from shutil import rmtree

def makepath(path):
    if not isdir(path):
        makedirs(path)


def removepath(path):
    if isdir(path):
        rmtree(path, ignore_errors=True)