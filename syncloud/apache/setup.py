import os
import platform

from syncloud.app import util
from syncloud.app.logger import get_logger
from syncloud.tools.apt import Apt
from ports import Ports
from env import sites_enabled_dir, http_file_template, http_file, http_web_root, \
    log_dir, http_include_dir, https_file_template, https_file, cert_file, key_file, \
    https_web_root


class Setup():
    def __init__(self, system):
        self.system = system
        self.logger = get_logger('system')
        self.apt = Apt()

    def install(self):

        (distname, version, id) = platform.linux_distribution()

        if distname == "Ubuntu" and version == "12.04":
            self.apt.add_repo("ppa:ondrej/php5-oldstable")
        self.apt.install(["libapache2-mod-wsgi", "php-apc", "curl", "libapache2-mod-php5"])

        self.system.enable_modules(["rewrite", "headers", "wsgi", "ssl"])
        self.system.init_conf_dirs()

    def activate(self, hostname):

        self.system.disable_all_sites(sites_enabled_dir)

        self.generate_config(http_file_template, http_file,
                             dict(
                                 hostname=hostname,
                                 web_root=http_web_root,
                                 log_dir=log_dir,
                                 include_dir=http_include_dir
                             ))
        self.system.enable_site("http")

        self.system.generate_certificate(hostname, cert_file, key_file)
        self.generate_config(https_file_template, https_file,
                             dict(
                                 hostname=hostname,
                                 cert_file=cert_file,
                                 key_file=key_file,
                                 web_root=https_web_root,
                                 log_dir=log_dir,
                                 include_dir=http_include_dir
                             ))
        self.system.enable_site("https")

        return Ports()

    def generate_config(self, from_filename, to_filename, mapping):

        self.logger.info('generate config: {}'.format(to_filename))
        util.transform_file(from_filename, to_filename, mapping)

        # Fix for apache 2.2
        to_filename_conf = "{}.conf".format(to_filename)
        if os.path.exists(to_filename_conf):
            os.remove(to_filename_conf)
        os.symlink(to_filename, to_filename_conf)