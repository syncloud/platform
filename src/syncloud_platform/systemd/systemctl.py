import os
from os.path import join
import shutil
from string import Template
import string
from subprocess import check_output, CalledProcessError
from syncloud_app import logger
from syncloud_platform.config.config import PlatformConfig

SYSTEMD_DIR = join('/lib', 'systemd', 'system')


def reload_service(service):

    log = logger.get_logger('systemctl')
    log.info('reloading {0}'.format(service))
    check_output('systemctl reload {0} 2>&1'.format(service), shell=True)


def remove_service(service):
    __remove('{0}.service'.format(service))


def __remove(filename):

    if "unknown" == __stop(filename):
        return

    check_output('systemctl disable {0} 2>&1'.format(filename), shell=True)
    systemd_file = __systemd_file(filename)
    if os.path.isfile(systemd_file):
        os.remove(systemd_file)


def add_service(app_id, service, include_socket=False, start=True):

    config = PlatformConfig()
    app_dir = join(config.apps_root(), app_id)

    log = logger.get_logger('systemctl')

    shutil.copyfile(__app_service_file(app_dir, service), __systemd_service_file(service))

    if include_socket:
        shutil.copyfile(__app_socket_file(app_dir, service), __systemd_socket_file(service))

    log.info('enabling {0}'.format(service))
    check_output('systemctl enable {0} 2>&1'.format(service), shell=True)
    if start:
        start_service(service)


def add_mount(mount_entry):

    log = logger.get_logger('systemctl')

    config = PlatformConfig()
    mount_template_file = join(config.config_dir(), 'mount', 'mount.template')
    mount_definition = Template(open(mount_template_file, 'r').read()).substitute({
        'what': mount_entry.device,
        'where': config.get_external_disk_dir(),
        'type': mount_entry.type,
        'options': mount_entry.options})

    config = PlatformConfig()
    mount_filename = __dir_to_systemd_mount_filename(config.get_external_disk_dir())
    with open(__systemd_file(mount_filename), 'w') as f:
        f.write(mount_definition)

    log.info('enabling {0}'.format(mount_filename))
    check_output('systemctl enable {0} 2>&1'.format(mount_filename), shell=True)
    __start(mount_filename)


def __dir_to_systemd_mount_filename(directory):
    return string.join(filter(None, directory.split('/')), '-') + '.mount'


def remove_mount():
    config = PlatformConfig()
    __remove(__dir_to_systemd_mount_filename(config.get_external_disk_dir()))


def restart_service(service):

    stop_service(service)
    start_service(service)


def start_service(service):
    __start('{0}.service'.format(service))


def start_mount(mount):
    __start('{0}.mount'.format(mount))


def __start(service):
    log = logger.get_logger('systemctl')

    try:
        log.info('starting {0}'.format(service))
        check_output('systemctl start {0} 2>&1'.format(service), shell=True)
    except CalledProcessError, e:
        try:
            log.error(check_output('systemctl status {0} 2>&1'.format(service), shell=True))
        except CalledProcessError, e:
            log.error(e.output)
        raise e


def stop_service(service):
    return __stop('{0}.service'.format(service))


def stop_mount(service):
    return __stop('{0}.mount'.format(service))


def __stop(service):
    log = logger.get_logger('systemctl')

    try:
        log.info('checking {0}'.format(service))
        result = check_output('systemctl is-active {0} 2>&1'.format(service), shell=True).strip()
        log.info('stopping {0}'.format(service))
        check_output('systemctl stop {0} 2>&1'.format(service), shell=True)
    except CalledProcessError, e:
        result = e.output.strip()

    log.info("{0}: {1}".format(service, result))
    return result


def __systemd_file(filename):
    return join(SYSTEMD_DIR, filename)


def __systemd_service_file(service):
    return __systemd_file("{0}.service".format(service))


def __systemd_socket_file(service):
    return join(SYSTEMD_DIR, "{0}.socket".format(service))


def __app_service_file(app_dir, service):
    return join(app_dir, 'config', 'systemd', "{0}.service".format(service))


def __app_socket_file(app_dir, service):
    return join(app_dir, 'config', 'systemd', "{0}.socket".format(service))
