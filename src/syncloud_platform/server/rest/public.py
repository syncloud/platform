import os
import traceback
from os.path import dirname, join, abspath
import sys

import convertible
from flask import Flask, jsonify, send_from_directory, request, redirect, send_file
from syncloud_platform.insider.config import InsiderConfig
from syncloud_platform.insider.facade import get_insider
from syncloud_platform.insider.redirect_service import RedirectService

from syncloud_platform.server.model import app_from_sam_app, App
from syncloud_platform.tools.hardware import Hardware

local_root = abspath(join(dirname(__file__), '..', '..', '..', '..'))
if __name__ == '__main__':
    sys.path.insert(0, local_root)

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.server.auth import authenticate
from syncloud_platform.sam.stub import SamStub
from syncloud_app import logger
from flask.ext.login import LoginManager, login_user, logout_user, current_user, login_required

sam = SamStub()

if __name__ == '__main__':
    www_dir = join(local_root, 'www', '_site')
    mock_apps = [
        App("owncloud", "ownCloud", "owncloud")]
    secret_key = '123223'
else:
    config = PlatformConfig()
    www_dir = config.www_root()
    mock_apps = None
    secret_key = config.get_web_secret_key()

html_prefix = '/server/html'
rest_prefix = '/server/rest'

logger.init(console=True)

app = Flask(__name__)
app.config['SECRET_KEY'] = secret_key
login_manager = LoginManager()
login_manager.init_app(app)


@login_manager.unauthorized_handler
def _callback():
    if request.is_xhr:
        return 'Unauthorised', 401
    else:
        return redirect(html_prefix + '/login.html')


class User:
    def __init__(self, name):
        self.name = name


class UserFlask:
    def __init__(self, user):
        self.user = user

    def is_authenticated(self):
        return True

    def is_active(self):
        return True

    def is_anonymous(self):
        return False

    def get_id(self):
        return unicode(self.user)


@app.route(html_prefix + '/<path:filename>')
def static_file(filename):
    return send_from_directory(www_dir, filename)


@login_manager.user_loader
def load_user(email):
    return UserFlask(User(email))


@app.route(rest_prefix + "/login", methods=["GET", "POST"])
def login():

    if not PlatformUserConfig().is_activated():
        return redirect('{0}://{1}:81'.format(request.scheme, request.host))

    if 'name' in request.form and 'password' in request.form:
        try:
            authenticate(request.form['name'], request.form['password'])
            user_flask = UserFlask(User(request.form['name']))
            login_user(user_flask, remember=False)
            # next_url = request.get('next_url', '/')
            return redirect("/")
        except Exception, e:
            traceback.print_exc(file=sys.stdout)
            return jsonify(message=e.message), 400

    return jsonify(message='missing name or password'), 400


@app.route(rest_prefix + "/logout", methods=["POST"])
@login_required
def logout():
    logout_user()
    return 'User logged out', 200


@app.route(rest_prefix + "/user", methods=["GET"])
@login_required
def user():
    return jsonify(convertible.to_dict(current_user.user)), 200


@app.route('/')
@login_required
def index():
    return static_file('index.html')

files_prefix = rest_prefix + '/files/'


@app.route(files_prefix)
@app.route(files_prefix + '<path:path>')
@login_required
def browser(path=''):
    filesystem_path = join('/', path)
    if os.path.isfile(filesystem_path):
        return send_file(filesystem_path, mimetype='text/plain')
    else:
        entries = os.listdir(filesystem_path)
        items = [{'name': entry, 'is_file': os.path.isfile(join(filesystem_path, entry))} for entry in entries]
        return jsonify(items=items, dir=filesystem_path)


@app.route(rest_prefix + "/installed_apps", methods=["GET"])
@login_required
def installed_apps():
    apps = [app_from_sam_app(a) for a in non_required_apps() if a.installed_version]

    # TODO: Hack to add system apps, need to think about it
    apps.append(App('store', 'App Store', html_prefix + '/store.html'))
    apps.append(App('settings', 'Settings', html_prefix + '/settings.html'))

    return jsonify(apps=convertible.to_dict(apps)), 200


@app.route(rest_prefix + "/app", methods=["GET"])
@login_required
def app_status():
    application = get_app(request.args['app_id'])
    return jsonify(info=convertible.to_dict(application)), 200


def get_app(app_id):
    return next(a for a in sam.list() if a.app.id == app_id)


@app.route(rest_prefix + "/install", methods=["GET"])
@login_required
def install():
    result = sam.install(request.args['app_id'])
    return jsonify(message=result), 200


@app.route(rest_prefix + "/remove", methods=["GET"])
@login_required
def remove():
    result = sam.remove(request.args['app_id'])
    return jsonify(message=result), 200


@app.route(rest_prefix + "/upgrade", methods=["GET"])
@login_required
def upgrade():
    result = sam.install(request.args['app_id'])
    return jsonify(message=result), 200


@app.route(rest_prefix + "/check", methods=["GET"])
@login_required
def update():
    result = sam.update()
    return jsonify(message=result), 200


@app.route(rest_prefix + "/available_apps", methods=["GET"])
@login_required
def available_apps():
    apps = [app_from_sam_app(a) for a in non_required_apps()]
    return jsonify(apps=convertible.to_dict(apps)), 200


@app.route(rest_prefix + "/settings/external_access", methods=["GET"])
@login_required
def get_settings_upnp():
    return jsonify(enabled=InsiderConfig().get_external_access()), 200


@app.route(rest_prefix + "/settings/external_access_enable", methods=["GET"])
@login_required
def external_access_enable():

    InsiderConfig().set_external_access(True)
    try:

        insider = get_insider()
        if not insider.mapper.available():
            return jsonify(success=False, message='No port mappers found (NatPmp, UPnP)'), 200

        insider.dns.sync()
        return jsonify(success=True), 200
    except Exception, e:
        InsiderConfig().set_external_access(False)
        return jsonify(success=False, message=e.message), 200


@app.route(rest_prefix + "/settings/external_access_disable", methods=["GET"])
@login_required
def external_access_disable():
    InsiderConfig().set_external_access(False)
    return jsonify(success=True), 200


@app.route(rest_prefix + "/send_log", methods=["GET"])
@login_required
def send_log():
    RedirectService().send_log()
    return jsonify(success=True), 200


@app.route(rest_prefix + "/settings/disks", methods=["GET"])
@login_required
def disks():
    return jsonify(success=True, disks=convertible.to_dict(Hardware().available_disks())), 200


@app.route(rest_prefix + "/settings/disk_activate", methods=["GET"])
@login_required
def disk_activate():
    device = request.args['device']
    fix_permissions = True
    if 'fix_permissions' in request.args:
        fix_permissions = request.args['fix_permissions'] == 'True'
    return jsonify(success=True, disks=Hardware().activate_disk(device, fix_permissions)), 200


@app.route(rest_prefix + "/settings/version", methods=["GET"])
@login_required
def version():
    return jsonify(convertible.to_dict(get_app('platform'))), 200


@app.route(rest_prefix + "/settings/system_upgrade", methods=["GET"])
@login_required
def system_upgrade():
    sam.upgrade('platform')
    return 'OK', 200


@app.route(rest_prefix + "/settings/disk_deactivate", methods=["GET"])
@login_required
def disk_deactivate():
    return jsonify(success=True, disks=Hardware().deactivate_disk()), 200


def non_required_apps():
    if mock_apps:
        apps = mock_apps
    else:
        apps = [a for a in sam.list() if not a.app.required]
    return apps


@app.errorhandler(Exception)
def handle_exception(error):
    print '-'*60
    traceback.print_exc(file=sys.stdout)
    print '-'*60
    response = jsonify(message=error.message)
    status_code = 500
    return response, status_code


def filter_websites(endpoints):
    return [endpoint for endpoint in endpoints
            if endpoint.service.protocol in ["http", "https"] and endpoint.service.name != "server"]


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
