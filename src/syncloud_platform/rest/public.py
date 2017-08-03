import sys
import traceback
from os import environ
import convertible
from flask import jsonify, send_from_directory, request, redirect, send_file, Flask
from flask.ext.login import login_user, logout_user, current_user, login_required
from flask_login import LoginManager

from syncloud_platform.auth.ldapauth import authenticate
from syncloud_platform.injector import get_injector
from syncloud_platform.rest.props import html_prefix, rest_prefix
from syncloud_platform.rest.flask_decorators import nocache, redirect_if_not_activated
from syncloud_platform.rest.model.flask_user import FlaskUser
from syncloud_platform.rest.model.user import User
from syncloud_platform.gaplib import linux

from syncloud_platform.rest.service_exception import ServiceException

injector = get_injector(environ['CONFIG_DIR'])
public = injector.public
device = injector.device

app = Flask(__name__)
app.config['SECRET_KEY'] = public.platform_config.get_web_secret_key()
login_manager = LoginManager()
login_manager.init_app(app)


@login_manager.unauthorized_handler
def _callback():
    if request.is_xhr:
        return 'Unauthorised', 401
    else:
        return redirect(html_prefix + '/login.html')


@app.route(html_prefix + '/<path:filename>')
@nocache
@redirect_if_not_activated
def static_file(filename):
    return send_from_directory(public.www_dir, filename)


@login_manager.user_loader
def load_user(email):
    return FlaskUser(User(email))


@app.route(rest_prefix + "/login", methods=["GET", "POST"])
@redirect_if_not_activated
def login():

    if 'name' in request.form and 'password' in request.form:
        try:
            authenticate(request.form['name'], request.form['password'])
            user_flask = FlaskUser(User(request.form['name']))
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


@app.route(rest_prefix + "/installed_apps", methods=["GET"])
@login_required
def installed_apps():
    return jsonify(apps=convertible.to_dict(public.installed_apps())), 200


@app.route(rest_prefix + "/app", methods=["GET"])
@login_required
def app_status():
    return jsonify(info=convertible.to_dict(public.get_app(request.args['app_id']))), 200


@app.route(rest_prefix + "/install", methods=["GET"])
@login_required
def install():
    public.install(request.args['app_id'])
    return 'OK', 200


@app.route(rest_prefix + "/remove", methods=["GET"])
@login_required
def remove():
    return jsonify(message=public.remove(request.args['app_id'])), 200


@app.route(rest_prefix + "/restart", methods=["GET"])
@login_required
def restart():
    public.restart()
    return 'OK', 200


@app.route(rest_prefix + "/shutdown", methods=["GET"])
@login_required
def shutdown():
    public.shutdown()
    return 'OK', 200


@app.route(rest_prefix + "/upgrade", methods=["GET"])
@login_required
def upgrade():
    public.upgrade(request.args['app_id'])
    return 'OK', 200


@app.route(rest_prefix + "/check", methods=["GET"])
@login_required
def update():
    return jsonify(message=public.update()), 200


@app.route(rest_prefix + "/available_apps", methods=["GET"])
@login_required
def available_apps():
    return jsonify(apps=convertible.to_dict(public.available_apps())), 200


@app.route(rest_prefix + "/access/access", methods=["GET"])
@login_required
def access():
    return jsonify(success=True, data=public.access()), 200


@app.route(rest_prefix + "/access/set_access", methods=["GET"])
@login_required
def set_access():
    public.set_access(
        request.args['upnp_enabled'] == 'true',
        request.args['is_https'] == 'true',
        request.args['external_access'] == 'true',
        request.args['public_ip'],
        int(request.args['public_port'])
    )
    return jsonify(success=True), 200


@app.route(rest_prefix + "/access/network_interfaces", methods=["GET"])
@login_required
def network_interfaces():
    return jsonify(success=True, data=dict(interfaces=public.network_interfaces())), 200


@app.route(rest_prefix + "/send_log", methods=["GET"])
@login_required
def send_log():
    include_support = request.args['include_support'] == 'true'
    public.send_logs(include_support)
    return jsonify(success=True), 200


@app.route(rest_prefix + "/settings/device_domain", methods=["GET"])
@login_required
def device_domain():
    return jsonify(success=True, device_domain=public.domain()), 200


@app.route(rest_prefix + "/settings/disks", methods=["GET"])
@login_required
def disks():
    return jsonify(success=True, disks=convertible.to_dict(public.disks())), 200


@app.route(rest_prefix + "/settings/boot_disk", methods=["GET"])
@login_required
def boot_disk():
    return jsonify(success=True, data=convertible.to_dict(public.boot_disk())), 200


@app.route(rest_prefix + "/settings/disk_activate", methods=["GET"])
@login_required
def disk_activate():
    return jsonify(success=True, disks=public.disk_activate(request.args['device'])), 200


@app.route(rest_prefix + "/settings/versions", methods=["GET"])
@login_required
def versions():
    return jsonify(success=True, data=convertible.to_dict(public.list_apps())), 200


@app.route(rest_prefix + "/settings/system_upgrade", methods=["GET"])
@login_required
def system_upgrade():
    public.system_upgrade()
    return 'OK', 200


@app.route(rest_prefix + "/settings/sam_upgrade", methods=["GET"])
@login_required
def sam_upgrade():
    public.sam_upgrade()
    return 'OK', 200


@app.route(rest_prefix + "/settings/sam_status", methods=["GET"])
@login_required
def sam_status():
    return jsonify(is_running=public.sam_status()), 200


@app.route(rest_prefix + "/settings/boot_extend", methods=["GET"])
@login_required
def boot_extend():
    return jsonify(is_running=public.boot_extend()), 200


@app.route(rest_prefix + "/settings/boot_extend_status", methods=["GET"])
@login_required
def boot_extend_status():
    return jsonify(is_running=public.boot_extend_status()), 200


@app.route(rest_prefix + "/settings/disk_deactivate", methods=["GET"])
@login_required
def disk_deactivate():
    return jsonify(success=True, disks=public.disk_deactivate()), 200


@app.route(rest_prefix + "/settings/regenerate_certificate", methods=["GET"])
@login_required
def regenerate_certificate():
    public.regenerate_certificate()
    return jsonify(success=True), 200


@app.route(rest_prefix + "/settings/activate_url", methods=["GET"])
@login_required
def activate_url():
    return jsonify(activate_url='{0}:81'.format(linux.local_ip()), success=True), 200


@app.errorhandler(Exception)
def handle_exception(error):
    print '-'*60
    traceback.print_exc(file=sys.stdout)
    print '-'*60
    status_code = 500

    if isinstance(error, ServiceException):
        status_code = 200

    response = jsonify(success=False, message=error.message)
    return response, status_code
