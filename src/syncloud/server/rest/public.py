import traceback
from os.path import dirname, join, abspath
import sys

import convertible
from flask import Flask, jsonify, send_from_directory, request, redirect

from syncloud.sam.manager import get_sam
from syncloud.sam.models import AppVersions
from syncloud.server.model import app_from_sam_app, App

local_root = abspath(join(dirname(__file__), '..', '..', '..'))
if __name__ == '__main__':
    sys.path.insert(0, local_root)

from syncloud.config.config import PlatformConfig
from syncloud.server.auth import authenticate
from syncloud.app import logger
from flask.ext.login import LoginManager, login_user, logout_user, current_user, login_required

if __name__ == '__main__':
    www_dir = join(local_root, 'www', '_site')
    mock_apps = [
        App("image-ci", "image ci", 'image-ci'),
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
login_manager.login_view = "/server/html/login.html"
login_manager.init_app(app)


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
    user = current_user.user
    user_data = convertible.to_dict(user)
    return jsonify(user_data), 200


@app.route('/')
@login_required
def index():
    return static_file('index.html')


@app.route(rest_prefix + "/installed_apps", methods=["GET"])
@login_required
def installed_apps():
    apps = [app_from_sam_app(app) for app in non_required_apps() if app.installed_version]

    # TODO: Hack to add system apps, need to think about it
    apps.append(App('store', 'App Store', 'server/html/store.html'))

    return jsonify(apps=convertible.to_dict(apps)), 200

@app.route(rest_prefix + "/available_apps", methods=["GET"])
@login_required
def available_apps():
    apps = [app_from_sam_app(app) for app in non_required_apps()]
    return jsonify(apps=convertible.to_dict(apps)), 200


def non_required_apps():
    if mock_apps:
        apps = mock_apps
    else:
        apps = get_sam().list()
        # TODO: pip-less sam should not prefix apps with 'syncloud-'
        for app in apps:
            app.app.id = app.app.id[9:]
        apps = [app for app in apps if not app.app.required]
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
    return [endpoint for endpoint in endpoints if endpoint.service.protocol in ["http", "https"] and endpoint.service.name != "server"]


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)