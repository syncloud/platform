import os
from os.path import join
import shutil
from string import Template
import string
from subprocess import check_output, CalledProcessError
from syncloud_app import logger

SYSTEMD_DIR = join('/lib', 'systemd', 'system')


class Systemctl:

    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.log = logger.get_logger('systemctl')

    def reload_service(self, service):
        service = self.service_name(service)
        
        self.log.info('reloading {0}'.format(service))
        try:
            check_output('systemctl reload {0} 2>&1'.format(service), shell=True)
        except CalledProcessError, e:
            self.log.error(e.output)
            raise e

    def remove_service(self, service):
        self.__remove('{0}.service'.format(service))

    def __remove(self, filename):

        if self.__stop(filename) in ("unknown", "inactive"):
            return

        check_output('systemctl disable {0} 2>&1'.format(filename), shell=True)
        systemd_file = self.__systemd_file(filename)
        if os.path.isfile(systemd_file):
            os.remove(systemd_file)

    def add_service(self, app_id, service, include_socket=False, start=True):

        service = self.service_name(service)

        configs_root = join(self.platform_config.configs_root(), app_id)

        log = logger.get_logger('systemctl')

        shutil.copyfile(self.__app_service_file(configs_root, service), self.__systemd_service_file(service))

        if include_socket:
            shutil.copyfile(self.__app_socket_file(configs_root, service), self.__systemd_socket_file(service))

        log.info('enabling {0}'.format(service))
        check_output('systemctl enable {0} 2>&1'.format(service), shell=True)
        if start:
            self.__start('{0}.service'.format(service))

    def add_mount(self, device, fs_type, options):

        log = logger.get_logger('systemctl')

        mount_template_file = join(self.platform_config.config_dir(), 'mount', 'mount.template')
        mount_definition = Template(open(mount_template_file, 'r').read()).substitute({
            'what': device,
            'where': self.platform_config.get_external_disk_dir(),
            # 'type': fs_type,
            'type': 'auto',
            'options': options})

        mount_filename = dir_to_systemd_mount_filename(self.platform_config.get_external_disk_dir())
        with open(self.__systemd_file(mount_filename), 'w') as f:
            f.write(mount_definition)

        log.info('enabling {0}'.format(mount_filename))
        check_output('systemctl enable {0} 2>&1'.format(mount_filename), shell=True)
        self.__start(mount_filename)

    def remove_mount(self, ):
        self.__remove(dir_to_systemd_mount_filename(self.platform_config.get_external_disk_dir()))

    def restart_service(self, service):

        self.stop_service(service)
        self.start_service(service)

    def start_service(self, service):
        service = self.service_name(service)
        self.__start('{0}.service'.format(service))

    def start_mount(self, mount):
        self.__start('{0}.mount'.format(mount))

    def __start(self, service):
        log = logger.get_logger('systemctl')

        try:
            log.info('starting {0}'.format(service))
            check_output('systemctl start {0} 2>&1'.format(service), shell=True)
        except CalledProcessError, e:
            log.error(e.output)
            try:
                log.error(check_output('journalctl -u {0} 2>&1'.format(service), shell=True))
            except CalledProcessError, e:
                log.error(e.output)
            raise e

    def stop_service(self, service):
        service = self.service_name(service)
        return self.__stop('{0}.service'.format(service))

    def stop_mount(self, service):
        return self.__stop('{0}.mount'.format(service))

    def __stop(self, service):
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

    def service_name(self, service):
        if self.platform_config.get_installer() == 'snapd':
            service = "snap.{0}".format(service)
        return service

    def __systemd_file(self, filename):
        return join(SYSTEMD_DIR, filename)

    def __systemd_service_file(self, service):
        return self.__systemd_file("{0}.service".format(service))

    def __systemd_socket_file(self, service):
        return join(SYSTEMD_DIR, "{0}.socket".format(service))

    def __app_service_file(self, app_dir, service):
        return join(app_dir, 'config', 'systemd', "{0}.service".format(service))

    def __app_socket_file(self, app_dir, service):
        return join(app_dir, 'config', 'systemd', "{0}.socket".format(service))


def dir_to_systemd_mount_filename(directory):
    return string.join(filter(None, directory.split('/')), '-') + '.mount'
