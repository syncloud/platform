import os
from syncloud.app.logger import get_logger
from syncloud.app import runner
from syncloud.apache.env import https_include_dir_full, http_include_dir_full, http_web_root, https_web_root


class System():

    def __init__(self):
        self.logger = get_logger('apache.system')

    def enable_module(self, name):
        runner.call("a2enmod {}".format(name), self.logger, shell=True)

    def enable_modules(self, names):
        for name in names:
            self.enable_module(name)

    def init_conf_dirs(self):
        self._create_dir(https_include_dir_full)
        self._create_dir(http_include_dir_full)
        self._create_dir(http_web_root)
        self._create_dir(https_web_root)

    def _create_dir(self, dir):
        if not os.path.exists(dir):
            self.logger.debug("creating: {0}".format(dir))
            os.mkdir(dir)

    def disable_all_sites(self, sites_enabled_dir):
        enabled_sites = os.listdir(sites_enabled_dir)
        if len(enabled_sites) > 0:
            self.logger.info('enabled sites: {}'.format(enabled_sites))
            for site in enabled_sites:
                runner.call("a2dissite {}".format(site), self.logger, shell=True)

    def enable_site(self, name):
        runner.call("a2ensite {}".format(name), self.logger, shell=True)

    def restart(self):
        self.logger.debug('restarting apache')
        runner.call("invoke-rc.d apache2 restart", self.logger, shell=True)

    def generate_certificate(self, hostname, cert_file, key_file):
        openssl_cmd = "openssl req -new -x509 -nodes"
        subj = "/C=US/ST=Unknown/L=Unknown/O=Unknown/CN={}".format(hostname)
        cmd = '{} -subj "{}" -out {} -keyout {}'.format(openssl_cmd, subj, cert_file, key_file)
        runner.call(cmd, self.logger, shell=True)