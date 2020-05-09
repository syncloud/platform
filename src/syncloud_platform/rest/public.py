import sys
import traceback

from syncloudlib.json import convertible
from syncloudlib.error import PassthroughJsonError

import requests
from flask import jsonify, request, redirect, Flask, Response
from flask_login import LoginManager, login_user, logout_user, current_user, login_required

from syncloud_platform.auth.ldapauth import authenticate
from syncloud_platform.injector import get_injector
from syncloud_platform.rest.flask_decorators import nocache, redirect_if_not_activated, redirect_if_activated
from syncloud_platform.rest.model.flask_user import FlaskUser
from syncloud_platform.rest.model.user import User
from syncloud_platform.gaplib import linux
from syncloud_platform.rest.backend_proxy import backend_request
from syncloud_platform.rest.service_exception import ServiceException
from syncloud_platform.rest.internal_validator import InternalValidator

injector = get_injector()
public = injector.public
internal = injector.internal
device = injector.device

app = Flask(__name__)
app.config['SECRET_KEY'] = public.user_platform_config.get_web_secret_key()
login_manager = LoginManager()
login_manager.init_app(app)


@login_manager.unauthorized_handler
def _callback():
    if request.is_xhr:
        return 'Unauthorised', 401
    else:
        return redirect('/login.html')


@app.route("/rest/id", methods=["GET"])
def identification():
    return jsonify(success=True, message='', data=convertible.to_dict(internal.identification())), 200


@app.route("/rest/activation_status", methods=["GET"])
def activation_status():
    try:
        return jsonify(activated=get_injector().user_platform_config.is_activated()), 200
    except Exception as e:
        return jsonify(activated=False), 200


@app.route("/rest/activate", methods=["POST"])
@redirect_if_activated
def activate():

    device_username = request.form['device_username'].lower()
    device_password = request.form['device_password']
        
    validator = InternalValidator()
    validator.validate(device_username, device_password)
    
    main_domain = 'syncloud.it'
    if 'main_domain' in request.form:
        main_domain = request.form['main_domain']

    internal.activate(
        request.form['redirect_email'],
        request.form['redirect_password'],
        request.form['user_domain'],
        device_username,
        device_password,
        main_domain
    )
    return identification()


@app.route("/rest/activate_custom_domain", methods=["POST"])
@redirect_if_activated
def activate_custom_domain():

    device_username = request.form['device_username'].lower()
    device_password = request.form['device_password']
        
    validator = InternalValidator()
    validator.validate(device_username, device_password)

    internal.activate_custom_domain(
        request.form['full_domain'],
        device_username,
        device_password,
    )
    return identification()


@login_manager.user_loader
def load_user(email):
    return FlaskUser(User(email))


@app.route("/rest/login", methods=["GET", "POST"])
@redirect_if_not_activated
def login():

    if 'name' in request.form and 'password' in request.form:
        try:
            authenticate(request.form['name'], request.form['password'])
            user_flask = FlaskUser(User(request.form['name']))
            login_user(user_flask, remember=False)
            # next_url = request.get('next_url', '/')
            return redirect("/")
        except Exception as e:
            traceback.print_exc(file=sys.stdout)
            return jsonify(message=str(e)), 400

    return jsonify(message='missing name or password'), 400


@app.route("/rest/logout", methods=["POST"])
@redirect_if_not_activated
@login_required
def logout():
    logout_user()
    return 'User logged out', 200


@app.route("/rest/user", methods=["GET"])
@redirect_if_not_activated
@login_required
def user():
    return jsonify(convertible.to_dict(current_user.user)), 200


@app.route("/rest/installed_apps", methods=["GET"])
@redirect_if_not_activated
@login_required
def installed_apps():
    return jsonify(apps=convertible.to_dict(public.installed_apps())), 200


@app.route("/rest/app", methods=["GET"])
@redirect_if_not_activated
@login_required
def app_status():
    return jsonify(info=convertible.to_dict(public.get_app(request.args['app_id']))), 200


@app.route("/rest/install", methods=["GET"])
@redirect_if_not_activated
@login_required
def install():
    public.install(request.args['app_id'])
    return 'OK', 200


@app.route("/rest/remove", methods=["GET"])
@redirect_if_not_activated
@login_required
def remove():
    return jsonify(message=public.remove(request.args['app_id'])), 200


@app.route("/rest/restart", methods=["GET"])
@redirect_if_not_activated
@login_required
def restart():
    public.restart()
    return 'OK', 200


@app.route("/rest/shutdown", methods=["GET"])
@redirect_if_not_activated
@login_required
def shutdown():
    public.shutdown()
    return 'OK', 200


@app.route("/rest/upgrade", methods=["GET"])
@redirect_if_not_activated
@login_required
def upgrade():

    force = False
    if 'force' in request.args:
        force = request.args['force'] == 'true'

    channel = public.platform_config.get_channel()
    if 'channel' in request.args:
        channel = request.args['channel']

    public.upgrade(request.args['app_id'], channel, force)

    return 'OK', 200


@app.route("/rest/available_apps", methods=["GET"])
@redirect_if_not_activated
@login_required
def available_apps():
    return jsonify(apps=convertible.to_dict(public.available_apps())), 200


@app.route("/rest/access/port_mappings", methods=["GET"])
@redirect_if_not_activated
@login_required
def port_mappings():
    return jsonify(success=True, port_mappings=convertible.to_dict(public.port_mappings())), 200


@app.route("/rest/access/access", methods=["GET"])
@redirect_if_not_activated
@login_required
def access():
    return jsonify(success=True, data=public.access()), 200


@app.route("/rest/access/set_access", methods=["GET"])
@redirect_if_not_activated
@login_required
def set_access():
    public_ip = None
    if 'public_ip' in request.args:
        public_ip = request.args['public_ip']
    public.set_access(
        request.args['upnp_enabled'] == 'true',
        request.args['external_access'] == 'true',
        public_ip,
        int(request.args['certificate_port']),
        int(request.args['access_port'])
    )
    return jsonify(success=True), 200


@app.route("/rest/access/network_interfaces", methods=["GET"])
@redirect_if_not_activated
@login_required
def network_interfaces():
    return jsonify(success=True, data=dict(interfaces=public.network_interfaces())), 200


@app.route("/rest/send_log", methods=["GET"])
@redirect_if_not_activated
@login_required
def send_log():
    include_support = request.args['include_support'] == 'true'
    public.send_logs(include_support)
    return jsonify(success=True), 200


@app.route("/rest/settings/device_domain", methods=["GET"])
@redirect_if_not_activated
@login_required
def device_domain():
    return jsonify(success=True, device_domain=public.domain()), 200


@app.route("/rest/settings/device_url", methods=["GET"])
@redirect_if_not_activated
@login_required
def device_url():
    return jsonify(success=True, device_url=public.device_url()), 200


@app.route("/rest/settings/disks", methods=["GET"])
@redirect_if_not_activated
@login_required
def disks():
    return jsonify(success=True, disks=convertible.to_dict(public.disks())), 200


@app.route("/rest/settings/boot_disk", methods=["GET"])
@redirect_if_not_activated
@login_required
def boot_disk():
    return jsonify(success=True, data=convertible.to_dict(public.boot_disk())), 200


@app.route("/rest/settings/disk_activate", methods=["GET"])
@redirect_if_not_activated
@login_required
def disk_activate():
    return jsonify(success=True, disks=public.disk_activate(request.args['device'])), 200


@app.route("/rest/settings/versions", methods=["GET"])
@redirect_if_not_activated
@login_required
def versions():
    return jsonify(success=True, data=convertible.to_dict(public.list_apps())), 200


@app.route("/rest/settings/installer_status", methods=["GET"])
@redirect_if_not_activated
@login_required
def installer_status():
    return jsonify(is_running=public.installer_status()), 200


@app.route("/rest/settings/disk_deactivate", methods=["GET"])
@redirect_if_not_activated
@login_required
def disk_deactivate():
    return jsonify(success=True, disks=public.disk_deactivate()), 200


@app.route("/rest/settings/regenerate_certificate", methods=["GET"])
@redirect_if_not_activated
@login_required
def regenerate_certificate():
    public.regenerate_certificate()
    return jsonify(success=True), 200


@app.route("/rest/settings/deactivate", methods=["POST"])
@redirect_if_not_activated
@login_required
def deactivate():
    public.user_platform_config.set_deactivated()
    return jsonify(success=True), 200


@app.route("/rest/app_image", methods=["GET"])
@redirect_if_not_activated
@login_required
def app_image():
    channel = request.args['channel']
    app = request.args['app']
    r = requests.get('http://apps.syncloud.org/releases/{0}/images/{1}-128.png'.format(channel, app), stream=True)
    return Response(r.iter_content(chunk_size=10*1024),
                    content_type=r.headers['Content-Type'])


@app.route("/rest/backup/<path:path>", methods=["GET"])
@app.route("/rest/installer/<path:path>", methods=["GET"])
@app.route("/rest/job/<path:path>", methods=["GET"])
@app.route("/rest/storage/<path:path>", methods=["POST"])
@redirect_if_not_activated
@login_required
def backend_proxy(path):
    response = backend_request(request.method, request.full_path.replace("/rest", "", 1), request.form)
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
