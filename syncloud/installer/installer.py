from subprocess import check_output

from syncloud.sam.installer import Installer
from syncloud.systemd.systemctl import add_service, remove_service
from syncloud.tools import app


class PlatformInstaller:

    def install(self, from_file=None):
        Installer().install('platform', from_file, 'platform')

        data_dir = app.get_app_data_root('platform', 'platform')
        app.create_data_dir(data_dir, 'webapps', 'platform')

        add_service('platform', 'platform-uwsgi-internal')
        add_service('platform', 'platform-uwsgi-public')
        add_service('platform', 'platform-nginx')
        add_service('platform', 'platform-openldap', start=False)

        check_output('syncloud-boot-installer', shell=True)

    def remove(self):

        remove_service('platform-openldap')
        remove_service('platform-nginx')
        remove_service('platform-uwsgi-public')
        remove_service('platform-uwsgi-internal')
