import logging
import shutil
from os.path import join, isdir
from subprocess import check_output, CalledProcessError

from syncloudlib import logger, fs

from syncloud_platform.config.user_config import PlatformUserConfig
from syncloud_platform.gaplib import linux, gen
from syncloud_platform.injector import get_injector

APP_NAME = 'platform'


class PlatformInstaller:
    def __init__(self):
        if not logger.factory_instance:
            logger.init(logging.DEBUG, True)

        self.log = logger.get_logger('installer')
        self.data_dir = '/var/snap/platform/current'
        self.common_dir = '/var/snap/platform/common'
        self.slapd_config_dir = join(self.data_dir, 'slapd.d')


    def init_configs(self):
        linux.fix_locale()

        data_dirs = [
            join(self.common_dir, 'log'),
            join(self.data_dir, 'nginx'),
            join(self.data_dir, 'openldap'),
            join(self.data_dir, 'openldap-data'),
            join(self.data_dir, 'certbot'),
            join(self.data_dir, 'certbot', 'www')
        ]

        for data_dir in data_dirs:
            fs.makepath(data_dir)
    
    def init_services(self):

        injector = get_injector()

        hardware = injector.hardware
        hardware.init_disk()

        injector.tls.init_certificate()

        ldap_auth = injector.ldap_auth
        ldap_auth.init()

        nginx = injector.nginx
        nginx.init_config()

    def install(self):
        self.init_configs()
        user_config = PlatformUserConfig()
        user_config.init_config()
        self.init_services()
        self.clear_crontab()

    def pre_refresh(self):
        self.clear_crontab()

    def clear_crontab(self):
        # crontab was migrated into backend process
        try:
            check_output('crontab -u root -r', shell=True)
        except CalledProcessError as e:
            self.log.error(e.output.decode())

    def post_refresh(self):
        self.init_configs()
        user_config = PlatformUserConfig()
        user_config.init_config()
        self.migrate_ldap()
        self.init_services()

    def migrate_ldap(self):
        old_config = '/var/snap/platform/common/slapd.d'
        if not isdir(self.slapd_config_dir):
            shutil.copytree(old_config, self.slapd_config_dir)
            # TODO: migrate module loading config
    def configure(self):
        pass
