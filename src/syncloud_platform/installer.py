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
            join(self.data_dir, 'openldap-data'),
            join(self.data_dir, 'certbot'),
            join(self.data_dir, 'certbot', 'www')
        ]

        for data_dir in data_dirs:
            fs.makepath(data_dir)

        check_output("/snap/platform/current/bin/update_certs.sh", shell=True)

    def init_services(self):

        injector = get_injector()

        hardware = injector.hardware
        hardware.init_disk()

        injector.tls.init_certificate()
        check_output("/snap/platform/current/bin/cli cert", shell=True)

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
        self.migrate_common_to_current()
        self.init_configs()
        user_config = PlatformUserConfig()
        user_config.init_config()
        self.init_services()

    def migrate_common_to_current(self):
        old_config_db = '/var/snap/platform/common/platform.db'
        new_config_db = '/var/snap/platform/current/platform.db'
        if not isfile(new_config_db):
            shutil.copyfile(old_config_db, new_config_db)

        old_certificate = '/var/snap/platform/common/syncloud.crt'
        new_certificate = '/var/snap/platform/current/syncloud.crt'
        if not isfile(new_certificate):
            shutil.copyfile(old_certificate, new_certificate)

        old_key = '/var/snap/platform/common/syncloud.key'
        new_key = '/var/snap/platform/current/syncloud.key'
        if not isfile(new_key):
            shutil.copyfile(old_key, new_key)

        old_ldap_data = '/var/snap/platform/common/openldap-data'
        new_ldap_data = '/var/snap/platform/current/openldap-data'
        if not isdir(new_ldap_data):
            shutil.copytree(old_ldap_data, new_ldap_data)

        old_slapd_config = '/var/snap/platform/common/slapd.d'
        if not isdir(self.slapd_config_dir):
            shutil.copytree(old_slapd_config, self.slapd_config_dir)

            shutil.copyfile(join(self.snap_dir, 'config/ldap/upgrade/cn=module{0}.ldif'),
                            join(self.slapd_config_dir, 'cn=config/cn=module{0}.ldif'))
            check_output('sed -i "s#{0}#{1}#g" {2}/cn=config.ldif'.format(self.common_dir, self.data_dir, self.slapd_config_dir), shell=True)
            check_output('sed -i "s#{0}#{1}#g" {2}/cn=config/olcDatabase={{1}}mdb.ldif'.format(self.common_dir, self.data_dir, self.slapd_config_dir), shell=True)

    def configure(self):
        pass
