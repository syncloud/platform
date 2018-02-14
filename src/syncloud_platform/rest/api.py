import sys
import traceback

from flask import Flask, jsonify, request, Response
from syncloud_app.main import PassthroughJsonError

from syncloud_platform.application.api import get_app_paths, get_app_setup
from syncloud_platform.injector import get_injector

app = Flask(__name__)


@app.route("/app/install_path", methods=["GET"])
def app_install_path():
    app_name = request.args['name']
    install_path = get_app_paths(app_name).get_install_dir()
    return jsonify(success=True, message='', data=install_path), 200


@app.route("/app/data_path", methods=["GET"])
def app_data_path():
    app_name = request.args['name']
    data_path = get_app_paths(app_name).get_data_dir()
    return jsonify(success=True, message='', data=data_path), 200


@app.route("/app/url", methods=["GET"])
def app_url():
    app_name = request.args['name']
    url = get_app_setup(app_name).app_url()
    return jsonify(success=True, message='', data=url), 200


@app.route("/app/domain_name", methods=["GET"])
def app_domain_name():
    app_name = request.args['name']
    domain_name = get_app_setup(app_name).app_domain_name()
    return jsonify(success=True, message='', data=domain_name), 200


@app.route("/app/init_storage", methods=["POST"])
def init_storage():
    app_name = request.form['app_name']
    user_name = request.form['user_name']
    app_storage_dir = get_app_setup(app_name).init_storage(user_name)
    return jsonify(success=True, message='', data=app_storage_dir), 200


@app.route("/app/storage_dir", methods=["GET"])
def storage_dir():
    app_name = request.args['name']
    app_storage_dir = get_app_setup(app_name).get_storage_dir()
    return jsonify(success=True, message='', data=app_storage_dir), 200


@app.route("/user/email", methods=["GET"])
def user_email():
    email = get_injector().user_platform_config.get_user_email()
    return jsonify(success=True, message='', data=email), 200


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
