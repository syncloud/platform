import sys
import traceback

from flask import Flask, jsonify, request, Response
from syncloudlib.error import PassthroughJsonError

from syncloud_platform.injector import get_injector

app = Flask(__name__)


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
        print('-'*60)
        traceback.print_exc(file=sys.stdout)
        print('-'*60)
        return jsonify(message=str(error)), status_code


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
