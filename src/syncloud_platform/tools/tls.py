import filecmp
from subprocess import check_output

from syncloud_platform.api import info
from syncloud_platform.config.config import PlatformConfig


class Tls:
    def __init__(self):
        self.platform_config = PlatformConfig()

    def generate_certificate(self):
        check_output('openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout {0} -out {1} -subj {2} 2>&1'.format(
                    self.platform_config.get_ssl_key_file(),
                    self.platform_config.get_ssl_certificate_file(),
                    "/C=US/ST=Syncloud/L=Syncloud/O=Syncloud/OU=Syncloud/CN=*.{0}".format(info.domain())
                ), shell=True)

    def is_default_certificate_installed(self):
        return filecmp.cmp(
            self.platform_config.get_ssl_certificate_file(),
            self.platform_config.get_default_ssl_certificate_file())
