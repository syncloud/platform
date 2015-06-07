from os.path import join, dirname
import traceback
import sys
import convertible

from flask import Flask, jsonify

if __name__ == '__main__':
    sys.path.insert(0, join(dirname(__file__), '..', '..', '..'))
from syncloud.tools.facade import Facade
from syncloud.app import logger

logger.init(console=True)

app = Flask(__name__)


@app.route("/id", methods=["GET"])
def id():
    return jsonify(convertible.to_dict(Facade().id())), 200

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
