import filecmp
import tempfile
from string import Template
from subprocess import check_output

from syncloud_app import util
from syncloud_app.logger import get_logger

from syncloud_platform.api import info
from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.tools.nginx import Nginx


class Tls:
    def __init__(self):
        self.log = get_logger('tls')
        self.platform_config = PlatformConfig()

    def generate_certificate(self):

        key_file = self.platform_config.get_ssl_key_file()
        check_output('openssl genrsa -out {0} 4096 2>&1'.format(key_file), shell=True)

        cert_file = self.platform_config.get_ssl_certificate_file()
        fd, temp_configfile = tempfile.mkstemp()
        util.transform_file(self.platform_config.get_openssl_config(), temp_configfile, {'domain': info.domain()})
        cmd = 'openssl req -new -x509 -days 3650 -config {0} -key {1} -out {2} 2>&1'.format(
            temp_configfile, key_file, cert_file)
        self.log.info('running: ' + cmd)
        check_output(cmd, shell=True)

        Nginx().reload()

    def is_default_certificate_installed(self):
        return filecmp.cmp(
            self.platform_config.get_ssl_certificate_file(),
            self.platform_config.get_default_ssl_certificate_file())
