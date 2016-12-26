import logging
import os
from os.path import isdir, join
from syncloud_app import logger
from syncloud_platform.gaplib import fs, linux, gen
import shutil
from syncloud_platform.injector import get_injector

logger.init(logging.DEBUG, console=True, line_format='%(message)s')
log = logger.get_logger('platform_post_install')

install_dir = os.environ['SNAP']
app_data = os.environ['SNAP_COMMON']
templates_path = join(install_dir, 'config.templates')
config_dir = join(app_data, 'config')

variables = {
    'apps_root': install_dir,
    'data_root': app_data,
    # not used in snap
    'configs_root': 'not_used',
    'config_root': app_data,
    'config_dir': config_dir,
    'app_dir': install_dir,
    'app_data': app_data
}
gen.generate_files(templates_path, config_dir, variables)

data_dirs = [
    join(app_data, 'webapps'),
    join(app_data, 'log'),
    join(app_data, 'nginx'),
    join(app_data, 'openldap'),
    join(app_data, 'openldap-data'),
    join(app_data, 'certbot'),
    join(app_data, 'slapd.d')
]

for data_dir in data_dirs:
    fs.makepath(data_dir)

injector = get_injector(config_dir=config_dir)

platform_config = injector.platform_config
hardware = injector.hardware
path_checker = injector.path_checker
ldap_auth = injector.ldap_auth

if not isdir(platform_config.get_disk_root()):
    os.mkdir(platform_config.get_disk_root())

if not isdir(platform_config.get_internal_disk_dir()):
    os.mkdir(platform_config.get_internal_disk_dir())

if not path_checker.external_disk_link_exists():
    hardware.relink_disk(
        platform_config.get_disk_link(),
        platform_config.get_internal_disk_dir())

if not os.path.exists(platform_config.get_ssl_certificate_file()):
    shutil.copy(platform_config.get_default_ssl_certificate_file(), platform_config.get_ssl_certificate_file())
    shutil.copy(platform_config.get_default_ssl_key_file(), platform_config.get_ssl_key_file())
