import sys
import traceback

from syncloudlib.json import convertible
import requests
from flask import jsonify, send_from_directory, request, redirect, Flask, Response
from flask_login import LoginManager, login_user, logout_user, current_user, login_required

from syncloud_platform.auth.ldapauth import authenticate
from syncloud_platform.injector import get_injector
from syncloud_platform.rest.props import html_prefix, rest_prefix
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
        return redirect(html_prefix + '/login.html')


@app.route(rest_prefix + "/id", methods=["GET"])
def identification():
    return jsonify(success=True, message='', data=convertible.to_dict(internal.identification())), 200


@app.route(rest_prefix + "/activate", methods=["POST"])
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


@app.route(rest_prefix + "/activate_custom_domain", methods=["POST"])
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


@app.route(html_prefix + '/activate.html')
@redirect_if_activated
@nocache
def activate_html(filename):
    return send_from_directory(public.www_dir, filename)


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

    force = False
    if 'force' in request.args:
        force = request.args['force'] == 'true'

    channel = public.platform_config.get_channel()
    if 'channel' in request.args:
        channel = request.args['channel']

    public.upgrade(request.args['app_id'], channel, force)

    return 'OK', 200


@app.route(rest_prefix + "/available_apps", methods=["GET"])
@login_required
def available_apps():
    return jsonify(apps=convertible.to_dict(public.available_apps())), 200


@app.route(rest_prefix + "/access/port_mappings", methods=["GET"])
@login_required
def port_mappings():
    return jsonify(success=True, port_mappings=convertible.to_dict(public.port_mappings())), 200


@app.route(rest_prefix + "/access/access", methods=["GET"])
@login_required
def access():
    return jsonify(success=True, data=public.access()), 200


@app.route(rest_prefix + "/access/set_access", methods=["GET"])
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


@app.route(rest_prefix + "/settings/device_url", methods=["GET"])
@login_required
def device_url():
    return jsonify(success=True, device_url=public.device_url()), 200


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


@app.route(rest_prefix + "/settings/installer_status", methods=["GET"])
@login_required
def installer_status():
    return jsonify(is_running=public.installer_status()), 200


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
    return jsonify(activate_url='http://{0}:81'.format(linux.local_ip()), success=True), 200


@app.route(rest_prefix + "/app_image", methods=["GET"])
@login_required
def app_image():
    channel = request.args['channel']
    app = request.args['app']
    r = requests.get('http://apps.syncloud.org/releases/{0}/images/{1}-128.png'.format(channel, app), stream=True)
    return Response(r.iter_content(chunk_size=10*1024),
                    content_type=r.headers['Content-Type'])


@app.route(rest_prefix + "/backup/<path:path>", methods=["GET"])
@app.route(rest_prefix + "/installer/<path:path>", methods=["GET"])
@app.route(rest_prefix + "/job/<path:path>", methods=["GET"])
@app.route(rest_prefix + "/storage/<path:path>", methods=["POST"])
@login_required
def backend_proxy(path):
    response = backend_request(request.method, request.full_path.replace(rest_prefix, "", 1), request.form)
    return response.text, response.status_code


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


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
