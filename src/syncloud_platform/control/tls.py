import filecmp
import tempfile
from subprocess import check_output

from syncloud_app import util
from syncloud_app.logger import get_logger


class Tls:
    def __init__(self, platform_config, info, nginx):
        self.info = info
        self.log = get_logger('tls')
        self.platform_config = platform_config
        self.nginx = nginx
        self.certbot_bin = '{0}/certbot/bin/certbot'.format(self.platform_config.app_dir())

    def generate_real_certificate(self):
        try:

            self.log.info('running certbot')
            output = check_output('{0} certonly --cert-path {1} --key-path {2} --webroot --webroot-path {3} -d {4}'.format(
                 self.certbot_bin,
                 self.platform_config.get_ssl_certificate_file(),
                 self.platform_config.get_ssl_key_file(),
                 self.platform_config.www_root(),
                 self.info.domain()), shell=True)
            self.log.info(output)

        except Exception, e:
            self.log.warn('unable to generate real certificate: {0}'.format(e.message))

    def generate_self_signed_certificate(self):

        key_file = self.platform_config.get_ssl_key_file()
        check_output('openssl genrsa -out {0} 4096 2>&1'.format(key_file), shell=True)

        cert_file = self.platform_config.get_ssl_certificate_file()
        fd, temp_configfile = tempfile.mkstemp()
        util.transform_file(self.platform_config.get_openssl_config(), temp_configfile, {'domain': self.info.domain()})
        cmd = 'openssl req -new -x509 -days 3650 -config {0} -key {1} -out {2} 2>&1'.format(
            temp_configfile, key_file, cert_file)
        self.log.info('running: ' + cmd)
        check_output(cmd, shell=True)

        self.nginx.reload()

    def is_default_certificate_installed(self):
        return filecmp.cmp(
            self.platform_config.get_ssl_certificate_file(),
            self.platform_config.get_default_ssl_certificate_file())
