import os
from os.path import join
import shutil
from subprocess import check_output, CalledProcessError
from syncloud.app import logger
from syncloud.config.config import PlatformConfig

SYSTEMD_DIR = join('/lib', 'systemd', 'system')


def reload_service(service):

    log = logger.get_logger('systemctl')
    log.info('reloading {0}'.format(service))
    check_output('systemctl reload {0}'.format(service), shell=True)


def remove_service(service):

    log = logger.get_logger('systemctl')

    if "unknown" == stop_service(service):
        return

    check_output('systemctl disable {0}'.format(service), shell=True)
    os.remove(__systemd_service_file(service))

def add_service(app_id, service, include_socket=False, start=True):

    config = PlatformConfig()
    app_dir = join(config.apps_root(), app_id)

    log = logger.get_logger('systemctl')

    shutil.copyfile(__app_service_file(app_dir, service), __systemd_service_file(service))

    if include_socket:
        shutil.copyfile(__app_socket_file(app_dir, service), __systemd_socket_file(service))

    log.info('enabling {0}'.format(service))
    check_output('systemctl enable {0}'.format(service), shell=True)
    if start:
        start_service(service)


def restart_service(service):

    log = logger.get_logger('systemctl')
    stop_service(service)
    start_service(service)


def start_service(service):
    log = logger.get_logger('systemctl')

    try:
        log.info('starting {0}'.format(service))
        check_output('systemctl start {0}'.format(service), shell=True)
    except CalledProcessError, e:
        try:
            log.error(check_output('systemctl status {0}'.format(service), shell=True))
        except CalledProcessError, e:
            log.error(e.output)
        raise e

def stop_service(service):
    log = logger.get_logger('systemctl')

    try:
        log.info('checking {0}'.format(service))
        result = check_output('systemctl is-active {0}'.format(service), shell=True).strip()
        log.info('stopping {0}'.format(service))
        check_output('systemctl stop {0}'.format(service), shell=True)
    except CalledProcessError, e:
        result = e.output.strip()

    log.info("{0}: {1}".format(service, result))
    return result

def __systemd_service_file(service):
    return join(SYSTEMD_DIR, "{0}.service".format(service))

def __systemd_socket_file(service):
    return join(SYSTEMD_DIR, "{0}.socket".format(service))

def __app_service_file(app_dir, service):
    return join(app_dir, 'config', 'systemd', "{0}.service".format(service))

def __app_socket_file(app_dir, service):
    return join(app_dir, 'config', 'systemd', "{0}.socket".format(service))
