import sys
import traceback

from flask import Flask, jsonify, request, Response
from syncloud_app.main import PassthroughJsonError

from syncloud_platform.application.api import get_app_paths, get_app_setup

app = Flask(__name__)


@app.route("/app/install_path", methods=["GET"])
def app_install_path():
    app_name = request.args['name']
    dir = get_app_paths(app_name).get_install_dir()
    return jsonify(success=True, message='', data=dir), 200


@app.route("/app/data_path", methods=["GET"])
def app_data_path():
    app_name = request.args['name']
    dir = get_app_paths(app_name).get_data_dir()
    return jsonify(success=True, message='', data=dir), 200


@app.route("/app/url", methods=["GET"])
def app_url():
    app_name = request.args['name']
    url = get_app_setup(app_name).app_url()
    return jsonify(success=True, message='', data=url), 200

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
