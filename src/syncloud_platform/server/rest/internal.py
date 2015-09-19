import logging
from os.path import join, dirname, abspath
import traceback
import sys
import convertible

from flask import Flask, jsonify, send_from_directory, request
from syncloud_app.main import PassthroughJsonError

local_root = abspath(join(dirname(__file__), '..', '..', '..'))
if __name__ == '__main__':
    sys.path.insert(0, local_root)

from syncloud_platform.tools.facade import Facade
from syncloud_app import logger
from syncloud_platform.config.config import PlatformConfig
from syncloud_platform.server.serverfacade import get_server

config = PlatformConfig()
if __name__ == '__main__':
    www_dir = join(local_root, 'www', '_site')
else:
    www_dir = config.www_root()

logger.init(filename=config.get_rest_internal_log())

app = Flask(__name__)


@app.route('/server/html/<path:filename>')
def static_file(filename):
    return send_from_directory(www_dir, filename)


@app.route("/server/rest/id", methods=["GET"])
def identification():
    return jsonify(success=True, message='', data=convertible.to_dict(Facade().id())), 200


@app.route("/server/rest/activate", methods=["POST"])
def activate():

    # TODO: validation

    api_url = None
    if 'api-url' in request.form:
        api_url = request.form['api-url']

    domain = None
    if 'domain' in request.form:
        domain = request.form['domain']

    get_server().activate(
        request.form['redirect-email'],
        request.form['redirect-password'],
        request.form['redirect-domain'],
        request.form['name'],
        request.form['password'],
        api_url,
        domain
    )
    return identification()


@app.errorhandler(Exception)
def handle_exception(error):
    response = None
    status_code = 500
    if isinstance(error, PassthroughJsonError):
        response = error.json
    else:
        print '-'*60
        traceback.print_exc(file=sys.stdout)
        print '-'*60
        response = jsonify(message=error.message)
    return response, status_code


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
