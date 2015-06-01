import traceback
import sys

from flask import Flask, jsonify

app = Flask(__name__)


@app.route("/ping", methods=["GET"])
def ping():
    return jsonify({"ping": "pong"}), 200

@app.errorhandler(Exception)
def handle_exception(error):
    print '-'*60
    traceback.print_exc(file=sys.stdout)
    print '-'*60
    response = jsonify(message=error.message)
    status_code = 500
    return response, status_code


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True, port=5001)
