import traceback
import convertible
from flask import Flask, jsonify, send_from_directory, request
from os.path import dirname, join
import sys
from syncloud.server.auth import Auth, authenticate

if __name__ == '__main__':
    sys.path.insert(0, join(dirname(__file__), '..', '..'))
from syncloud.insider.dns import Endpoint
from syncloud.server.model import Site
from syncloud.app import logger
from syncloud.insider.config import Service, os
from syncloud.insider.facade import get_insider
from flask.ext.login import LoginManager, login_user, logout_user, current_user, login_required

www_dir = '/var/www/syncloud-server'
mock_sites = None
html_prefix = ''
rest_prefix = ''

if __name__ == '__main__':
    html_prefix = '/server/html'
    rest_prefix = '/server/rest'

logger.init(console=True)

app = Flask(__name__)
app.config['SECRET_KEY'] = '123223'
login_manager = LoginManager()
login_manager.login_view = html_prefix + "/login.html"
login_manager.init_app(app)


class User():
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
            authenticate(request.form['name'], request.form['password'], get_insider().full_name())
            user_flask = UserFlask(User(request.form['name']))
            login_user(user_flask, remember=False)
            return 'User logged in', 200
        except Exception, e:
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


@app.route(rest_prefix + "/sites", methods=["GET"])
@login_required
def site_list():
    if mock_sites:
        endpoints = mock_sites
    else:
        endpoints = filter_websites(get_insider().endpoints())

    return jsonify(sites=convertible.to_dict(map(Site, endpoints))), 200


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
    www_dir = join(dirname(__file__), '..', '..', 'www')
    mock_sites = [
        Endpoint(Service("image-ci", "http", "type", "80", "image-ci"), 'localhost', 8181),
        Endpoint(Service("owncloud", "https", "type", "443", "owncloud"), 'localhost', 8282)]
    app.run(host='0.0.0.0', debug=True, port=5001)