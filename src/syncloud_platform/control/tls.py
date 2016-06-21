import filecmp
import os
import subprocess
import tempfile
from subprocess import check_output, CalledProcessError

from os.path import join
from syncloud_app import util
from syncloud_app.logger import get_logger


def apps_to_certbot_domain_args(app_versions, domain):
    # we need to list all the individual domains for now as wildcard domain is not supported by certbot yet
    all_apps = [app_versions.app.id for app_versions in app_versions]
    domains = ['{0}.{1}'.format(app, domain) for app in all_apps]
    domains.append(domain)
    domains.reverse()
    domain_args = '-d ' + ' -d '.join(domains)
    return domain_args


class Tls:
    def __init__(self, platform_config, user_platform_config, info, nginx, sam):
        self.info = info
        self.log = get_logger('tls')
        self.platform_config = platform_config
        self.user_platform_config = user_platform_config
        self.nginx = nginx
        self.openssl_bin = '{0}/openssl/bin/openssl'.format(self.platform_config.app_dir())
        self.certbot_bin = '{0}/bin/certbot'.format(self.platform_config.app_dir())
        self.log_dir = self.platform_config.get_log_root()
        self.certbot_config_dir = join(self.platform_config.data_dir(), 'certbot')
        self.sam = sam

    def generate_real_certificate(self):

        if not self.platform_config.is_certbot_enabled():
            self.log.info('certbot is not enabled, not running')
            return

        if not self.user_platform_config.get_external_access():
            self.log.info('external access is not enabled, not running certbot')
            return

        try:

            self.log.info('running certbot')
            domain_args = apps_to_certbot_domain_args(self.sam.list(), self.info.domain())
            output = check_output(
                '{0} --logs-dir={1} --config-dir={2} --agree-tos --email {3} '
                'certonly --webroot --webroot-path {6} '
                '{7} '.format(self.certbot_bin,
                              self.log_dir,
                              self.certbot_config_dir,
                              self.user_platform_config.get_user_email(),
                              self.platform_config.www_root(),
                              domain_args), stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)

            if 'no action taken' not in output:
                certbot_certificate_file = '{0}/certbot/live/{1}/fullchain.pem'.format(
                    self.platform_config.data_dir(), self.info.domain())
                if os.path.exists(self.platform_config.get_ssl_certificate_file()):
                    os.remove(self.platform_config.get_ssl_certificate_file())
                os.symlink(certbot_certificate_file, self.platform_config.get_ssl_certificate_file())

                certbot_key_file = '{0}/certbot/keys/0000_key-certbot.pem'.format(
                    self.platform_config.data_dir(), self.info.domain())
                if os.path.exists(self.platform_config.get_ssl_key_file()):
                    os.remove(self.platform_config.get_ssl_key_file())
                os.symlink(certbot_key_file, self.platform_config.get_ssl_key_file())

                self.nginx.reload()

        except CalledProcessError, e:
            self.log.warn('unable to generate real certificate: {0}'.format(e))
            self.log.warn(e.output)

    def generate_self_signed_certificate(self):

        key_file = self.platform_config.get_ssl_key_file()
        try:

            output = check_output('{0} genrsa -out {1} 4096 2>&1'.format(self.openssl_bin, key_file),
                                  stderr=subprocess.STDOUT, shell=True)
            self.log.info(output)
        except CalledProcessError, e:
            self.log.warn('unable to generate self-signed certificate: {0}'.format(e))
            self.log.warn(e.output)
            raise e

        cert_file = self.platform_config.get_ssl_certificate_file()
        fd, temp_configfile = tempfile.mkstemp()
        util.transform_file(self.platform_config.get_openssl_config(), temp_configfile, {'domain': self.info.domain()})
        cmd = '{0} req -new -x509 -days 3650 -config {1} -key {2} -out {3} 2>&1'.format(self.openssl_bin,
                                                                                        temp_configfile, key_file,
                                                                                        cert_file)
        self.log.info('running: ' + cmd)
        output = check_output(cmd, shell=True)
        self.log.info(output)

        self.nginx.reload()

    def is_default_certificate_installed(self):
        return filecmp.cmp(
            self.platform_config.get_ssl_certificate_file(),
            self.platform_config.get_default_ssl_certificate_file())
