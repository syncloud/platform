from syncloud.app import runner
from syncloud.app import logger

from util import copyfile, create_link, rmfile, Settings

from os import makedirs
from os.path import exists, join
from shutil import rmtree

from keygen import KeyGen

from psutil import process_iter

BITS = 2048


class SshInitD:
    def __init__(self, filename):
        self.log = logger.get_logger('SshInitD')
        self.filename = filename
        copyfile('/etc/init.d/ssh', self.filename, self.log)
        self.__load()
        self.__replace('/etc/init.d/ssh', self.filename)

    def __load(self):
        with open(self.filename, 'r') as f:
            self.data = f.read()

    def __replace(self, old, new):
        self.data = self.data.replace(old, new)

    def pid_file(self, filename):
        self.__replace('/var/run/sshd.pid', filename)

    def running_folder(self, path):
        self.__replace('/var/run/sshd', path)

    def executable(self, filename):
        self.__replace('/usr/sbin/sshd', filename)

    def service_name(self, name):
        self.__replace(' sshd ', ' {0} '.format(name))
        self.__replace('"sshd"', '"{0}"'.format(name))
        self.__replace('"sshd"', '"{0}"'.format(name))
        self.__replace('Provides:		sshd', 'Provides:		{0}'.format(name))

    def default_config(self, filename):
        self.__replace('/etc/default/ssh', filename)

    def save(self):
        with open(self.filename, 'w') as f:
            f.write(self.data)


def update_defaults(name, log):
    runner.call('update-rc.d -f {0} defaults'.format(name), log, shell=True)


def service_remove(name, log):
    runner.call('insserv --verbose --remove {0}'.format(name), log, shell=True)


def service_add(name, log):
    runner.call('insserv --verbose --default {0}'.format(name), log, shell=True)


def add_to_initd(name, log):
    runner.call('insserv --verbose --remove {0}'.format(name), log, shell=True)


def services_reload(log):
    runner.call('systemctl --system daemon-reload', log, shell=True)


def service_start(name, log):
    runner.call('service {0} start'.format(name), log, shell=True)


def service_stop(name, log):
    runner.call('service {0} stop'.format(name), log, shell=True)


def restart_service(name, log):
    service_stop(name, log)
    service_start(name, log)


def yes_no(value):
    if value:
        return 'yes'
    else:
        return 'no'


class SshServer:
    def __init__(self, name):
        self.log = logger.get_logger('SshServer')
        self.key_gen = KeyGen()

        self.ssh_name = 'ssh_{0}'.format(name)
        self.sshd_name = 'sshd_{0}'.format(name)

        self.sshd_bin_filename = '/usr/sbin/{0}'.format(self.sshd_name)
        self.sshd_pam_filename = '/etc/pam.d/{0}'.format(self.sshd_name)

        self.sshd_pid_filename = '/var/run/sshd_{0}.pid'.format(name)

        self.ssh_config_filename = '/etc/ssh/{0}_config'.format(self.ssh_name)
        self.sshd_config_filename = '/etc/ssh/{0}_config'.format(self.sshd_name)

        self.certificates_path = '/root/.ssh_{0}'.format(name)
        self.authorized_keys_filename = join(self.certificates_path, 'authorized_keys')

        self.sshd_working_folder = '/var/run/{0}'.format(self.sshd_name)

        self.ssh_initd_filename = '/etc/init.d/{0}'.format(self.ssh_name)
        self.ssh_default_filename = '/etc/default/{0}'.format(self.ssh_name)

    def get_server_keys(self):
        rsa_key = '/etc/ssh/ssh_syncloud_host_rsa_key'
        dsa_key = '/etc/ssh/ssh_syncloud_host_dsa_key'
        ecdsa_key = '/etc/ssh/ssh_syncloud_host_ecdsa_key'
        if not exists(rsa_key):
            self.key_gen.generate_into_file('rsa', rsa_key, BITS, overwrite=True)
        if not exists(dsa_key):
            self.key_gen.generate_into_file('dsa', dsa_key, overwrite=True)
        if not exists(ecdsa_key):
            self.key_gen.generate_into_file('ecdsa', ecdsa_key, overwrite=True)
        return rsa_key, dsa_key, ecdsa_key

    def setup(self, service_port, password_authentication):
        rsa_key, dsa_key, ecdsa_key = self.get_server_keys()

        copyfile('/etc/ssh/ssh_config', self.ssh_config_filename, self.log)

        copyfile('/etc/ssh/sshd_config', self.sshd_config_filename, self.log)
        sshd_config = Settings(self.sshd_config_filename, delimeter=' ')
        sshd_config.set('PasswordAuthentication', yes_no(password_authentication))
        sshd_config.set('PubkeyAuthentication', yes_no(not password_authentication))
        if not password_authentication:
            sshd_config.set('AuthorizedKeysFile', self.authorized_keys_filename)
        sshd_config.set('Port', service_port)
        sshd_config.set('PidFile', self.sshd_pid_filename)
        sshd_config.set('HostKey', rsa_key, add_to_existing=False)
        sshd_config.set('HostKey', dsa_key, add_to_existing=True)
        sshd_config.set('HostKey', ecdsa_key, add_to_existing=True)
        sshd_config.save()

        create_link('/usr/sbin/sshd', self.sshd_bin_filename, self.log)
        create_link('/etc/pam.d/sshd', self.sshd_pam_filename, self.log)

        copyfile('/etc/default/ssh', self.ssh_default_filename, self.log)
        default_parameter = Settings(self.ssh_default_filename, delimeter='=')
        default_parameter.set('SSHD_OPTS', '"-f {0}"'.format(self.sshd_config_filename))
        default_parameter.save()
        update_defaults(self.ssh_name, self.log)

        if not exists(self.sshd_working_folder):
            makedirs(self.sshd_working_folder)

        ssh_init_d = SshInitD(self.ssh_initd_filename)
        ssh_init_d.pid_file(self.sshd_pid_filename)
        ssh_init_d.running_folder(self.sshd_working_folder)
        ssh_init_d.executable(self.sshd_bin_filename)
        ssh_init_d.service_name(self.sshd_name)
        ssh_init_d.default_config(self.ssh_default_filename)
        ssh_init_d.save()

        service_remove(self.ssh_name, self.log)
        service_add(self.ssh_name, self.log)
        services_reload(self.log)

        service_stop(self.ssh_name, self.log)
        service_start(self.ssh_name, self.log)

        if not password_authentication:
            return self.add_certificate()

    def add_certificate(self):
        if not exists(self.certificates_path):
            makedirs(self.certificates_path)

        private, public = self.key_gen.generate('rsa', BITS)

        with open(self.authorized_keys_filename, 'a+') as keys_file:
            keys_file.write(public)
            keys_file.write('\n')

        private_key_filename = join(self.certificates_path, 'id_rsa_syncloud_master')
        with open(private_key_filename, 'w+') as private_key_file:
            private_key_file.write(private)

        return private

    def remove(self):
        service_stop(self.ssh_name, self.log)

        sshd_proc = next((p for p in process_iter() if p.name() == self.sshd_name), None)
        if sshd_proc:
            sshd_proc.kill()

        service_remove(self.ssh_name, self.log)

        rmfile(self.ssh_initd_filename, ignore_errors=True)
        rmfile(self.ssh_default_filename, ignore_errors=True)

        rmfile(self.sshd_bin_filename, ignore_errors=True)
        rmfile(self.sshd_pam_filename, ignore_errors=True)

        rmfile(self.ssh_config_filename, ignore_errors=True)
        rmfile(self.sshd_config_filename, ignore_errors=True)

        rmtree(self.certificates_path, ignore_errors=True)
