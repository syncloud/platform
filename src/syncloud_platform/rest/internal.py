import traceback
import sys
from syncloudlib.json import convertible
from flask import Flask, jsonify, send_from_directory, request, Response
from syncloudlib.error import PassthroughJsonError

from syncloud_platform.injector import get_injector
from syncloud_platform.rest.props import rest_prefix, html_prefix
from syncloud_platform.rest.flask_decorators import nocache
from syncloud_platform.rest.internal_validator import InternalValidator

injector = get_injector()
internal = injector.internal
device = injector.device

app = Flask(__name__)


@app.route(html_prefix + "/<path:filename>")
@nocache
def static_file(filename):
    return send_from_directory(internal.www_dir, filename)


@app.route(rest_prefix + "/id", methods=["GET"])
def identification():
    return jsonify(success=True, message='', data=convertible.to_dict(internal.identification())), 200


@app.route(rest_prefix + "/activate", methods=["POST"])
def activate():

    device_username = request.form['device_username'].lower()
    device_password = request.form['device_password']
        
    validator = InternalValidator()
    validator.validate(device_username, device_password)
    
    main_domain = get_main_domain(request.form)

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


def get_main_domain(request_form):
    
    main_domain = None
    if 'main_domain' in request_form:
        if request_form['main_domain']:
            main_domain = request_form['main_domain']

    if main_domain is None:
        main_domain = "syncloud.it"

    return main_domain


@app.route(rest_prefix + "/send_log", methods=["POST"])
def send_log():
    internal.send_logs(
        request.form['redirect_email'],
        request.form['redirect_password'],
        get_main_domain(request.form))

    return jsonify(success=True), 200


@app.errorhandler(Exception)
def handle_exception(error):
    status_code = 500
    if isinstance(error, PassthroughJsonError):
        return Response(error.json, status=status_code, mimetype='application/json')
    else:
        print '-'*60
        traceback.print_exc(file=sys.stdout)
        print '-'*60
        return jsonify(message=error.message), status_code


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
