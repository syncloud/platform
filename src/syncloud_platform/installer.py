import logging
from os.path import join
from subprocess import check_output, CalledProcessError

from syncloudlib import logger, fs

from syncloud_platform.config import config
from syncloud_platform.config.user_config import PlatformUserConfig
from syncloud_platform.gaplib import linux, gen
from syncloud_platform.injector import get_injector

APP_NAME = 'platform'


class PlatformInstaller:
    def __init__(self):
        if not logger.factory_instance:
            logger.init(logging.DEBUG, True)

        self.log = logger.get_logger('installer')
        self.templates_path = join(config.INSTALL_DIR, 'config.templates')
        self.config_dir = join(config.DATA_DIR, 'config')
        self.data_dir = config.DATA_DIR
    
    def init_configs(self):
        linux.fix_locale()
        
        variables = {
            'apps_root': config.APPS_ROOT,
            'data_root': config.DATA_ROOT,
            'configs_root': config.DATA_ROOT,
            'config_root': self.data_dir,
            'config_dir': self.config_dir,
            'app_dir': config.INSTALL_DIR,
            'app_data': config.DATA_DIR
        }
        gen.generate_files(self.templates_path, self.config_dir, variables)

        data_dirs = [
            join(self.data_dir, 'webapps'),
            join(self.data_dir, 'log'),
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
        self.init_services()

    def configure(self):
        pass
