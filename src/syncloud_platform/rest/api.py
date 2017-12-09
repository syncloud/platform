import traceback
import sys
from os import environ
import convertible
from flask import Flask, jsonify, send_from_directory, request, Response
from syncloud_app.main import PassthroughJsonError

from syncloud_platform.application.api import get_app_paths
from syncloud_platform.rest.flask_decorators import nocache


app = Flask(__name__)


@app.route("/app/paths", methods=["GET"])
def app_paths():
    app_name = request.args['name']
    paths = get_app_paths(app_name)
    return jsonify(success=True, message='', data=paths), 200


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
