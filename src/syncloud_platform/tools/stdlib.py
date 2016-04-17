from os import makedirs, chown, utime
from os.path import isdir
from grp import getgrnam
from pwd import getpwnam
from shutil import rmtree

def makepath(path):
    if not isdir(path):
        makedirs(path)


def removepath(path):
    if isdir(path):
        rmtree(path, ignore_errors=True)


def chownpath(path, user):
    chown(path, getpwnam(user).pw_uid, getgrnam(user).gr_gid)


def touch(file, user):
    with open(file, 'a'):
        utime(file, None)
    chownpath(file, user)
