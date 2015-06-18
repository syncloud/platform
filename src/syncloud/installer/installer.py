from syncloud.sam.installer import Installer
from syncloud.systemd.systemctl import add_service, remove_service
from syncloud.tools import app


class PlatformInstaller:

    def install(self, from_file=None):
        Installer().install('platform', from_file, 'platform')

        data_dir = app.get_app_data_root('platform', 'platform')
        app.create_data_dir(data_dir, 'webapps', 'platform')

        add_service('platform', 'platform-cpu-frequency')
        add_service('platform', 'platform-insider-sync')
        add_service('platform', 'platform-resize-sd', start=False)
        add_service('platform', 'platform-udisks-glue')
        add_service('platform', 'platform-ntpdate')
        add_service('platform', 'platform-uwsgi-internal')
        add_service('platform', 'platform-uwsgi-public')
        add_service('platform', 'platform-nginx')
        add_service('platform', 'platform-openldap', start=False)

        # TODO: External disk mount
        # config = PlatformConfig()
        # LOCAL_DATA=/home/www-data
        # mkdir -p ${LOCAL_DATA}
        # chown -R www-data. ${LOCAL_DATA}
        # /usr/local/bin/syncloud-link-data ${LOCAL_DATA}
        # chmod 640 /etc/sudoers.d/www-data
        # adduser www-data plugdev
        # '/etc/sudoers.d', ['config/sudoers.d/www-data']),
        # '/etc/polkit-1/localauthority/50-local.d', ['config/polkit/55-storage.pkla']),
        # '/etc/udev/rules.d', ['config/udev/99-syncloud.udisks.rules']),

    def remove(self):

        remove_service('platform-openldap')
        remove_service('platform-nginx')
        remove_service('platform-uwsgi-public')
        remove_service('platform-uwsgi-internal')
        remove_service('platform-ntpdate')
        remove_service('platform-udisks-glue')
        remove_service('platform-resize-sd')
        remove_service('platform-insider-sync')
        remove_service('platform-cpu-frequency')
