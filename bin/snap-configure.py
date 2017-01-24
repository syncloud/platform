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
    join(app_data, 'certbot')
]

for data_dir in data_dirs:
    fs.makepath(data_dir)

injector = get_injector(config_dir=config_dir)

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
