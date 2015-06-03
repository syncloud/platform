from subprocess import check_output
from syncloud.insider.facade import get_insider
from syncloud.remote.remoteaccess import RemoteAccess

from syncloud.sam.installer import Installer
from syncloud.systemd.systemctl import add_service
from syncloud.app import logger

APP_ROOT = '/opt/app/platform'


class PlatformInstaller:

    def install(self, from_file=None):
        log = logger.get_logger('platform.postinstall')

        Installer().install('platform', from_file)

        add_service(APP_ROOT, 'platform-uwsgi-internal')
        add_service(APP_ROOT, 'platform-uwsgi-public')
        add_service(APP_ROOT, 'platform-nginx')

        check_output('syncloud_ssh_keys_generate', shell=True)

        check_output('syncloud-boot-installer', shell=True)

    def remove(self):
        RemoteAccess(get_insider()).disable()

