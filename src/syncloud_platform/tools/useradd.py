import pwd
from subprocess import check_output


def useradd(user):
    try:
        pwd.getpwnam(user)
        return 'user {0} exists'.format(user)
    except KeyError:
        return check_output('/usr/sbin/useradd -r -s /bin/false {0}'.format(user), shell=True)
