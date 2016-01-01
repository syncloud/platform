import os
import traceback
from os.path import dirname, join, abspath
import sys

import convertible
from flask import Flask, jsonify, send_from_directory, request, redirect, send_file
from syncloud_platform.insider.facade import get_insider
from syncloud_platform.insider.redirect_service import RedirectService

from syncloud_platform.rest.model import app_from_sam_app, App
from syncloud_platform.rest.flask_decorators import nocache, redirect_if_not_activated
from syncloud_platform.tools.hardware import Hardware
from syncloud_platform.tools.tls import Tls

local_root = abspath(join(dirname(__file__), '..', '..', '..', '..'))
if __name__ == '__main__':
    sys.path.insert(0, local_root)

from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig
from syncloud_platform.auth.ldapauth import authenticate
from syncloud_platform.sam.stub import SamStub
from syncloud_app import logger
from flask.ext.login import LoginManager, login_user, logout_user, current_user, login_required

config = PlatformConfig()
if __name__ == '__main__':
    www_dir = join(local_root, 'www', '_site')
    mock_apps = [
        App("owncloud", "ownCloud", "owncloud")]
    secret_key = '123223'
else:
    www_dir = config.www_root()
    mock_apps = None
    secret_key = config.get_web_secret_key()

html_prefix = '/server/html'
rest_prefix = '/server/rest'

logger.init(filename=config.get_rest_public_log())
log = logger.get_logger('rest.public')

app = Flask(__name__)
app.config['SECRET_KEY'] = secret_key
login_manager = LoginManager()
login_manager.init_app(app)
sam = SamStub()


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
@nocache
@redirect_if_not_activated
def static_file(filename):
    return send_from_directory(www_dir, filename)


@login_manager.user_loader
def load_user(email):
    return UserFlask(User(email))


@app.route(rest_prefix + "/login", methods=["GET", "POST"])
@redirect_if_not_activated
def login():

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
@redirect_if_not_activated
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
        entries = sorted(os.listdir(filesystem_path))
        items = [{'name': entry, 'is_file': os.path.isfile(join(filesystem_path, entry))} for entry in entries]
        return jsonify(items=items, dir=filesystem_path)


@app.route(rest_prefix + "/installed_apps", methods=["GET"])
@login_required
def installed_apps():
    apps = [app_from_sam_app(a) for a in sam.installed_user_apps()]

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
    return sam.get_app(app_id)


@app.route(rest_prefix + "/install", methods=["GET"])
@login_required
def install():
    sam.install(request.args['app_id'])
    return 'OK', 200


@app.route(rest_prefix + "/remove", methods=["GET"])
@login_required
def remove():
    result = sam.remove(request.args['app_id'])
    return jsonify(message=result), 200


@app.route(rest_prefix + "/upgrade", methods=["GET"])
@login_required
def upgrade():
    sam.upgrade(request.args['app_id'])
    return 'OK', 200


@app.route(rest_prefix + "/check", methods=["GET"])
@login_required
def update():
    result = sam.update()
    return jsonify(message=result), 200


@app.route(rest_prefix + "/available_apps", methods=["GET"])
@login_required
def available_apps():
    return jsonify(apps=convertible.to_dict([app_from_sam_app(a) for a in sam.user_apps()])), 200


@app.route(rest_prefix + "/settings/external_access", methods=["GET"])
@login_required
def external_access():
    return jsonify(external_access=PlatformUserConfig().get_external_access()), 200


@app.route(rest_prefix + "/settings/set_external_access", methods=["GET"])
@login_required
def external_access_enable():
    get_insider().add_main_device_service(PlatformUserConfig().get_protocol(), request.args['external_access'])
    return jsonify(success=True), 200


@app.route(rest_prefix + "/settings/protocol", methods=["GET"])
@login_required
def protocol():
    return jsonify(protocol=PlatformUserConfig().get_protocol()), 200


@app.route(rest_prefix + "/settings/set_protocol", methods=["GET"])
@login_required
def set_protocol():
    get_insider().add_main_device_service(request.args['protocol'], PlatformUserConfig().get_external_access())
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
    return jsonify(success=True, disks=Hardware().activate_disk(device)), 200


@app.route(rest_prefix + "/settings/version", methods=["GET"])
@login_required
def version():
    return jsonify(convertible.to_dict(get_app('platform'))), 200


@app.route(rest_prefix + "/settings/system_upgrade", methods=["GET"])
@login_required
def system_upgrade():
    sam.upgrade('platform')
    return 'OK', 200


@app.route(rest_prefix + "/settings/sam_status", methods=["GET"])
@login_required
def sam_status():
    return jsonify(is_running=sam.is_running()), 200


@app.route(rest_prefix + "/settings/disk_deactivate", methods=["GET"])
@login_required
def disk_deactivate():
    return jsonify(success=True, disks=Hardware().deactivate_disk()), 200


@app.errorhandler(Exception)
def handle_exception(error):
    print '-'*60
    traceback.print_exc(file=sys.stdout)
    print '-'*60
    response = jsonify(success=False, message=error.message)
    status_code = 500
    return response, status_code


# def filter_websites(endpoints):
#     return [endpoint for endpoint in endpoints
#             if endpoint.service.protocol in ["http", "https"] and endpoint.service.name != "server"]


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
