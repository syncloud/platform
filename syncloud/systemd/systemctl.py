import os
from os.path import join
import shutil
from subprocess import check_output, CalledProcessError

SYSTEMD_DIR = join('/lib', 'systemd', 'system')

def remove_service(service):

    try:
        check_output('systemctl is-active {0}'.format(service), shell=True)
        check_output('systemctl stop {0}'.format(service), shell=True)
    except CalledProcessError, e:
        result = e.output.strip()
        print("{0}: {1}".format(service, result))
        if result == "unknown":
            return

    os.remove(__systemd_service_file(service))

def add_service(app_dir, service):
    shutil.copyfile(__app_service_file(app_dir, service), __systemd_service_file(service))
    check_output('systemctl enable -f {0}'.format(service), shell=True)
    check_output('systemctl daemon-reload', shell=True)
    check_output('systemctl start {0}'.format(service), shell=True)

def __systemd_service_file(service):
    return join(SYSTEMD_DIR, "{0}.service".format(service))

def __app_service_file(app_dir, service):
    return join(app_dir, 'config', service)
