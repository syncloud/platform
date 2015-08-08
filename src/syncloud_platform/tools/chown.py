from subprocess import check_output
from syncloud_platform.tools.useradd import useradd


def chown(user, dir):
    useradd(user)
    return check_output('chown -RL {0}. {1}'.format(user, dir), shell=True)
