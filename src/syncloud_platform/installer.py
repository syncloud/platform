import logging
import shutil
from os.path import join, isfile, isdir
from subprocess import check_output, CalledProcessError

from syncloudlib import logger, fs

from syncloud_platform.config.user_config import PlatformUserConfig
from syncloud_platform.injector import get_injector

APP_NAME = 'platform'


class PlatformInstaller:
    def __init__(self):
        if not logger.factory_instance:
            logger.init(logging.DEBUG, True)

        self.log = logger.get_logger('installer')
        self.snap_dir = '/snap/platform/current'
        self.data_dir = '/var/snap/platform/current'
        self.common_dir = '/var/snap/platform/common'
        self.slapd_config_dir = join(self.data_dir, 'slapd.d')

    def init_configs(self):
        data_dirs = [
            join(self.common_dir, 'log'),
            join(self.data_dir, 'nginx'),
            join(self.data_dir, 'openldap'),
            join(self.data_dir, 'openldap-data')
        ]

        for data_dir in data_dirs:
            fs.makepath(data_dir)

        check_output("/snap/platform/current/bin/update_certs.sh", shell=True)

    def install(self):
        self.init_configs()
        user_config = PlatformUserConfig()
        user_config.init_config()
        injector = get_injector()
        injector.hardware.init_disk()
        check_output("/snap/platform/current/bin/cli cert", shell=True)
        injector.ldap_auth.init()
        injector.nginx.init_config()

    def pre_refresh(self):
        pass

    def post_refresh(self):
        self.init_configs()
        user_config = PlatformUserConfig()
        user_config.init_config()
        injector = get_injector()
        injector.hardware.init_disk()
        injector.nginx.init_config()

    
    def configure(self):
        pass
