import sys
import traceback

from syncloudlib.json import convertible
from syncloudlib.error import PassthroughJsonError

import requests
from flask import jsonify, request, redirect, Flask, Response
from flask_login import LoginManager, login_user, logout_user, current_user, login_required

from syncloud_platform.injector import get_injector
from syncloud_platform.rest.flask_decorators import nocache, fail_if_not_activated, fail_if_activated
from syncloud_platform.rest.model.flask_user import FlaskUser
from syncloud_platform.rest.model.user import User
from syncloud_platform.rest.backend_proxy import backend_request
from syncloud_platform.rest.service_exception import ServiceException
from syncloudlib.logger import get_logger

injector = get_injector()
public = injector.public


app = Flask(__name__)
app.config['SECRET_KEY'] = public.user_platform_config.get_web_secret_key()
login_manager = LoginManager()
login_manager.init_app(app)
log = get_logger('ldap')


@login_manager.unauthorized_handler
def _callback():
    log.warn('Unauthorised handler 401')
    return 'Unauthorised', 401


@login_manager.user_loader
def load_user(user_id):
    log.info('loading user {0}'.format(user_id))
    return FlaskUser(User(user))


@app.route("/rest/login", methods=["POST"])
@fail_if_not_activated
def login():
    request_json = request.json
    if 'username' in request_json and 'password' in request_json:
        try:
            injector.ldap_auth.authenticate(request_json['username'], request_json['password'])
            user_flask = FlaskUser(User(request_json['username']))
            log.info('login user {0}'.format(user_flask.user.name))
            login_user(user_flask, remember=False)
            # next_url = request.get('next_url', '/')
            return redirect("/")
        except Exception as e:
            traceback.print_exc(file=sys.stdout)
            return jsonify(message=str(e)), 400

    return jsonify(message='missing name or password'), 400


@app.route("/rest/logout", methods=["POST"])
@fail_if_not_activated
@login_required
def logout():
    log.info('logout user {0}'.format(current_user.user.name))
    logout_user()
    return 'User logged out', 200


@app.route("/rest/user", methods=["GET"])
@fail_if_not_activated
@login_required
def user():
    log.info('current user {0}'.format(current_user.user.name))
    return jsonify(convertible.to_dict(current_user.user)), 200


@app.route("/rest/app", methods=["GET"])
@fail_if_not_activated
@login_required
def app_status():
    return jsonify(info=convertible.to_dict(public.get_app(request.args['app_id']))), 200


@app.route("/rest/install", methods=["POST"])
@fail_if_not_activated
@login_required
def install():
    public.install(request.json['app_id'])
    return jsonify(success=True), 200


@app.route("/rest/remove", methods=["POST"])
@fail_if_not_activated
@login_required
def remove():
    return jsonify(message=public.remove(request.json['app_id'])), 200


@app.route("/rest/restart", methods=["POST"])
@fail_if_not_activated
@login_required
def restart():
    public.restart()
    return jsonify(success=True), 200


@app.route("/rest/shutdown", methods=["POST"])
@fail_if_not_activated
@login_required
def shutdown():
    public.shutdown()
    return jsonify(success=True), 200


@app.route("/rest/upgrade", methods=["POST"])
@fail_if_not_activated
@login_required
def upgrade():
    public.upgrade(request.json['app_id'])
    return jsonify(success=True), 200


@app.route("/rest/access/network_interfaces", methods=["GET"])
@fail_if_not_activated
@login_required
def network_interfaces():
    return jsonify(success=True, data=dict(interfaces=public.network_interfaces())), 200


@app.route("/rest/send_log", methods=["POST"])
@fail_if_not_activated
@login_required
def send_log():
    include_support = request.args['include_support'] == 'true'
    public.send_logs(include_support)
    return jsonify(success=True), 200


@app.route("/rest/settings/device_domain", methods=["GET"])
@fail_if_not_activated
@login_required
def device_domain():
    return jsonify(success=True, device_domain=public.domain()), 200


@app.route("/rest/settings/device_url", methods=["GET"])
@fail_if_not_activated
@login_required
def device_url():
    return jsonify(success=True, device_url=public.device_url()), 200


@app.route("/rest/settings/boot_disk", methods=["GET"])
@fail_if_not_activated
@login_required
def boot_disk():
    return jsonify(success=True, data=convertible.to_dict(public.boot_disk())), 200


@app.route("/rest/settings/disk_activate", methods=["POST"])
@fail_if_not_activated
@login_required
def disk_activate():
    return jsonify(success=True, disks=public.disk_activate(request.json['device'])), 200


@app.route("/rest/settings/installer_status", methods=["GET"])
@fail_if_not_activated
@login_required
def installer_status():
    return jsonify(is_running=public.installer_status()), 200


@app.route("/rest/settings/disk_deactivate", methods=["POST"])
@fail_if_not_activated
@login_required
def disk_deactivate():
    return jsonify(success=True, disks=public.disk_deactivate()), 200


@app.route("/rest/settings/deactivate", methods=["POST"])
@fail_if_not_activated
@login_required
def deactivate():
    log.info('deactivate')
    logout_user()
    public.user_platform_config.set_deactivated()
    return jsonify(success=True), 200


@app.route("/rest/app_image", methods=["GET"])
@fail_if_not_activated
@login_required
def app_image():
    channel = request.args['channel']
    app = request.args['app']
    r = requests.get('http://apps.syncloud.org/releases/{0}/images/{1}-128.png'.format(channel, app), stream=True)
    return Response(r.iter_content(chunk_size=10*1024),
                    content_type=r.headers['Content-Type'])


@app.route("/rest/backup/list", methods=["GET"])
@app.route("/rest/backup/create", methods=["POST"])
@app.route("/rest/backup/restore", methods=["POST"])
@app.route("/rest/backup/remove", methods=["POST"])
@app.route("/rest/installer/upgrade", methods=["POST"])
@app.route("/rest/installer/version", methods=["GET"])
@app.route("/rest/job/status", methods=["GET"])
@app.route("/rest/storage/disk_format", methods=["POST"])
@app.route("/rest/storage/boot_extend", methods=["POST"])
@app.route("/rest/storage/disks", methods=["GET"])
@app.route("/rest/event/trigger", methods=["POST"])
@app.route("/rest/certificate", methods=["GET"])
@app.route("/rest/certificate/log", methods=["GET"])
@app.route("/rest/access", methods=["GET", "POST"])
@app.route("/rest/apps/available", methods=["GET"])
@app.route("/rest/apps/installed", methods=["GET"])
@fail_if_not_activated
@login_required
def backend_proxy_activated():
    response = backend_request(request.method, request.full_path.replace("/rest", "", 1), request.json)
    return response.text, response.status_code


@app.route("/rest/redirect/domain/availability", methods=["POST"])
@app.route("/rest/redirect_info", methods=["GET"])
@app.route("/rest/activate/managed", methods=["POST"])
@app.route("/rest/activate/custom", methods=["POST"])
@fail_if_activated
def backend_proxy_not_activated():
    response = backend_request(request.method, request.full_path.replace("/rest", "", 1), request.json)
    return response.text, response.status_code


@app.route("/rest/activation/status", methods=["GET"])
def backend_proxy():
    response = backend_request(request.method, request.full_path.replace("/rest", "", 1), request.json)
    return response.text, response.status_code


@app.route("/rest/id", methods=["GET"])
def identification():
    response = backend_request("GET", "/id", None)
    return response.text, response.status_code


@app.errorhandler(Exception)
def handle_exception(error):
    print('-'*60)
    traceback.print_exc(file=sys.stdout)
    print('-'*60)
    status_code = 500

    if isinstance(error, PassthroughJsonError):
        return Response(error.json, status=status_code, mimetype='application/json')

    if isinstance(error, ServiceException):
        status_code = 200

    response = jsonify(success=False, message=str(error))
    return response, status_code


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
