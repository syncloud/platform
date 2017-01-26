import filecmp
import os
import shutil
import subprocess
import tempfile
from subprocess import check_output, CalledProcessError

from syncloud_app import util
from syncloud_app.logger import get_logger
import time


class Tls:
    def __init__(self, platform_config, user_platform_config, info, nginx, certbot_generator):
        self.info = info
        self.log = get_logger('tls')
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.nginx = nginx
        self.openssl_bin = '{0}/openssl/bin/openssl'.format(self.platform_config.app_dir())
        self.certbot_generator = certbot_generator

    def generate_real_certificate(self):

        try:
            self.log.info('running certbot')
            result = self.certbot_generator.generate_certificate(self.platform_config.is_certbot_test_cert())

            if result.regenerated:
                if os.path.exists(self.platform_config.get_ssl_certificate_file()):
                    os.remove(self.platform_config.get_ssl_certificate_file())
                os.symlink(result.certificate_file, self.platform_config.get_ssl_certificate_file())

                if os.path.exists(self.platform_config.get_ssl_key_file()):
                    os.remove(self.platform_config.get_ssl_key_file())
                os.symlink(result.key_file, self.platform_config.get_ssl_key_file())

                self.nginx.reload()

        except CalledProcessError, e:
            self.log.warn('unable to generate real certificate: {0}'.format(e))
            self.log.warn(e.output)

    def generate_self_signed_certificate(self):

        try:

            key_ca_file = self.platform_config.get_ssl_ca_key_file()
            cert_ca_file = self.platform_config.get_ssl_ca_certificate_file()
            ssl_ca_serial_file = self.platform_config.get_ssl_ca_serial_file()

            with open(ssl_ca_serial_file, 'w') as f:
                f.write(str(int(time.time())))

            fd, temp_configfile = tempfile.mkstemp()
            util.transform_file(self.platform_config.get_openssl_config(), temp_configfile,
                                {
                                    'domain': self.info.domain(),
                                    'config_dir': self.platform_config.config_dir(),
                                    'ssl_ca_key_file': key_ca_file,
                                    'ssl_ca_certificate_file': cert_ca_file
                                })

            self.log.info('generating CA Key')
            output = check_output('OPENSSL_CONF={2} {0} genrsa -out {1} 4096 2>&1'.format(
                self.openssl_bin,
                key_ca_file,
                temp_configfile),
                stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)

            self.log.info('generating CA Certificate')
            output = check_output('OPENSSL_CONF={1} {0} req -new -x509 -days 3650 -config {1} -key {2} -out {3} 2>&1'
                                  .format(self.openssl_bin,
                                          temp_configfile,
                                          key_ca_file,
                                          cert_ca_file),
                                  stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)

            self.log.info('generating Server Key')
            key_file = self.platform_config.get_ssl_key_file()
            output = check_output('OPENSSL_CONF={2} {0} genrsa -out {1} 4096 2>&1'
                                  .format(self.openssl_bin,
                                          key_file,
                                          temp_configfile),
                                  stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)

            self.log.info('generating Server Certificate Request')
            key_file = self.platform_config.get_ssl_key_file()
            certificate_request_file = self.platform_config.get_ssl_certificate_request_file()
            output = check_output('OPENSSL_CONF={1} {0} req -config {1} -key {2} -new -sha256 -out {3} 2>&1'
                                  .format(self.openssl_bin,
                                          temp_configfile,
                                          key_file,
                                          certificate_request_file),
                                  stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)

            self.log.info('generating Server Certificate')
            cert_file = self.platform_config.get_ssl_certificate_file()
            output = check_output(
                'OPENSSL_CONF={1} {0} ca -config {1} '
                '-extensions server_cert -days 3650 '
                '-notext -md sha256 -in {2} -out {3} -batch 2>&1'.format(
                    self.openssl_bin,
                    temp_configfile,
                    certificate_request_file,
                    cert_file),
                stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)

        except CalledProcessError, e:
            self.log.warn('unable to generate self-signed certificate: {0}'.format(e))
            self.log.warn(e.output)
            raise e

        self.nginx.reload()

    def is_default_certificate_installed(self):
        return filecmp.cmp(
            self.platform_config.get_ssl_certificate_file(),
            self.platform_config.get_default_ssl_certificate_file())

    def init_certificate(self):
        if not os.path.exists(self.platform_config.get_ssl_certificate_file()):
            shutil.copy(self.platform_config.get_default_ssl_certificate_file(), self.platform_config.get_ssl_certificate_file())
            shutil.copy(self.platform_config.get_default_ssl_key_file(), self.platform_config.get_ssl_key_file())