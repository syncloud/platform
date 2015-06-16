from sys import exec_prefix
from os.path import join


def usr_local_dir():
    if exec_prefix == '/usr':
        return join(exec_prefix, 'local')
    else:
        return exec_prefix


def root_dir_prefix():
    if exec_prefix == '/usr':
        return ''
    else:
        return exec_prefix
