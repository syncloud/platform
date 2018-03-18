import os
from subprocess import check_output
from os.path import isdir, join
import shutil
from syncloud_app import logger
from syncloud_platform.injector import get_injector
from syncloud_platform.application import api

from syncloud_platform.gaplib import fs, linux, gen

# this makes config file paths relative

APP_NAME = 'platform'


class PlatformInstaller:
    def __init__(self):
        self.log = logger.get_logger('platform_post_install')

    def install(self):
        linux.fix_locale()
        
        if 'SNAP' in os.environ:
            apps_root = '/snap'
            data_root = '/var/snap'
            install_dir = os.environ['SNAP']
            data_dir = os.environ['SNAP_COMMON']
            app_data_prefix = 'common/'
        else:
            apps_root = '/opt/app'
            data_root = '/opt/data'
            install_dir = join(apps_root, APP_NAME)
            data_dir = join(data_root, APP_NAME)
            app_data_prefix = ''
            
        templates_path = join(install_dir, 'config.templates')
        config_dir = join(data_dir, 'config')

        variables = {
            'apps_root': apps_root,
            'data_root': data_root,
            'configs_root': data_root,
            'config_root': data_dir,
            'config_dir': config_dir,
            'app_dir': install_dir,
            'app_data': data_dir,
            'app_data_prefix': app_data_prefix
        }
        gen.generate_files(templates_path, config_dir, variables)

        data_dirs = [
            join(data_dir, 'webapps'),
            join(data_dir, 'log'),
            join(data_dir, 'nginx'),
            join(data_dir, 'openldap'),
            join(data_dir, 'openldap-data'),
            join(data_dir, 'certbot')
        ]

        for data_dir in data_dirs:
            fs.makepath(data_dir)

        injector = get_injector()

        hardware = injector.hardware
        hardware.init_disk()

        injector.tls.init_certificate()

        injector.platform_cron.remove()
        injector.platform_cron.create()

        udev = injector.udev
        udev.remove()
        udev.add()

        ldap_auth = injector.ldap_auth
        ldap_auth.init()

        nginx = injector.nginx
        nginx.init_config()
        
    def start(self):
        injector = get_injector()

        systemctl = injector.systemctl
        systemctl.add_service(APP_NAME, 'platform.cpu-frequency')
        systemctl.add_service(APP_NAME, 'platform.insider-sync')
        systemctl.add_service(APP_NAME, 'platform.ntpdate')
        systemctl.add_service(APP_NAME, 'platform.uwsgi-api')
        systemctl.add_service(APP_NAME, 'platform.uwsgi-internal')
        systemctl.add_service(APP_NAME, 'platform.uwsgi-public')
        systemctl.add_service(APP_NAME, 'platform.nginx-api')
        systemctl.add_service(APP_NAME, 'platform.nginx-internal')
        systemctl.add_service(APP_NAME, 'platform.nginx-public')
        systemctl.add_service(APP_NAME, 'platform.openldap')

    def remove(self):
        injector = get_injector()
        systemctl = injector.systemctl

        systemctl.remove_service('platform.openldap')
        systemctl.remove_service('platform.nginx-public')
        systemctl.remove_service('platform.nginx-internal')
        systemctl.remove_service('platform.nginx-api')
        systemctl.remove_service('platform.uwsgi-public')
        systemctl.remove_service('platform.uwsgi-internal')
        systemctl.remove_service('platform.uwsgi-api')
        systemctl.remove_service('platform.ntpdate')
        systemctl.remove_service('platform.insider-sync')
        systemctl.remove_service('platform.cpu-frequency')

        injector.platform_cron.remove()
        injector.udev.remove()

