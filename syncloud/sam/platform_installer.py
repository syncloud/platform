from subprocess import check_output

from syncloud.insider.facade import get_insider
from syncloud.remote.remoteaccess import RemoteAccess
from syncloud.sam.installer import Installer
from syncloud.systemd.systemctl import add_service
from syncloud.tools import app


class PlatformInstaller:

    def install(self, from_file=None):
        Installer().install('platform', from_file)

        data_dir = app.get_app_data_root('platform')
        app.create_data_dir(data_dir, 'webapps')

        add_service('platform', 'platform-uwsgi-internal')
        add_service('platform', 'platform-uwsgi-public')
        add_service('platform', 'platform-nginx')

        check_output('syncloud_ssh_keys_generate', shell=True)

        check_output('syncloud-boot-installer', shell=True)

    def remove(self):
        RemoteAccess(get_insider()).disable()

