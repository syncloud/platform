import os

from os.path import join

from flask import Flask
from flask_login import LoginManager
from syncloud_app import logger

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig, PLATFORM_APP_NAME
from syncloud_platform.device import get_device
from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.insider.service_config import ServiceConfig
from syncloud_platform.rest.model.app import app_from_sam_app
from syncloud_platform.sam.models import App
from syncloud_platform.sam.stub import SamStub
from syncloud_platform.tools import network
from syncloud_platform.tools.app import get_app_data_root
from syncloud_platform.tools.hardware import Hardware

html_prefix = '/server/html'
rest_prefix = '/server/rest'


class Public:

    def __init__(self):
        self.platform_config = PlatformConfig()
        logger.init(filename=self.platform_config.get_rest_public_log())
        self.log = logger.get_logger('rest.public')
        self.user_platform_config = PlatformUserConfig(self.platform_config.get_user_config())
        self.device = get_device()
        self.sam = SamStub()

        self.flask_app = Flask(__name__)
        self.flask_app.config['SECRET_KEY'] = self.platform_config.get_web_secret_key()
        self.flask_login_manager = LoginManager()
        self.flask_login_manager.init_app(self.flask_app)

        self.www_dir = self.platform_config.www_root()

    def browse(self, filesystem_path):
        entries = sorted(os.listdir(filesystem_path))
        return [{'name': entry, 'is_file': os.path.isfile(join(filesystem_path, entry))} for entry in entries]

    def installed_apps(self):
        apps = [app_from_sam_app(a) for a in self.sam.installed_user_apps()]

        # TODO: Hack to add system apps, need to think about it
        apps.append(App('store', 'App Store', html_prefix + '/store.html'))
        apps.append(App('settings', 'Settings', html_prefix + '/settings.html'))
        return apps

    def get_app(self, app_id):
        return self.sam.get_app(app_id)

    def install(self, app_id):
        self.sam.install(app_id)

    def remove(self, app_id):
        return self.sam.remove(app_id)

    def upgrade(self, app_id):
        self.sam.upgrade(app_id)

    def update(self):
        return self.sam.update()

    def available_apps(self):
        return [app_from_sam_app(a) for a in self.sam.user_apps()]

    def external_access(self):
        return self.user_platform_config.get_external_access()

    def external_access_enable(self, external_access):
        self.device.set_access(self.user_platform_config.get_protocol(), external_access)

    def protocol(self):
        return self.user_platform_config.get_protocol()

    def set_protocol(self, protocol):
        self.device.set_access(protocol, self.user_platform_config.get_external_access())

    def send_log(self):
        data_root = get_app_data_root(PLATFORM_APP_NAME)
        service_config = ServiceConfig(data_root)
        redirect_service = RedirectService(service_config, network.local_ip(), self.user_platform_config, self.platform_config)
        get_user_update_token = self.user_platform_config.get_user_update_token()
        redirect_service.send_log(get_user_update_token)

    def disk_activate(self, device):
        return Hardware().activate_disk(device)

    def system_upgrade(self):
        self.sam.upgrade('platform')

    def sam_status(self):
        return self.sam.is_running()

    def disk_deactivate(self):
        return Hardware().deactivate_disk()
