from sys import exec_prefix
from os.path import join


def is_virtual_env():
    return not exec_prefix == '/usr'


def usr_local_dir():
    if not is_virtual_env():
        return join(exec_prefix, 'local')
    else:
        return exec_prefix


def root_dir_prefix():
    if not is_virtual_env():
        return ''
    else:
        return exec_prefix
