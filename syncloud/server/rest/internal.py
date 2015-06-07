from os.path import join, dirname, abspath
import traceback
import sys
import convertible

from flask import Flask, jsonify, send_from_directory

local_root = abspath(join(dirname(__file__), '..', '..', '..'))
if __name__ == '__main__':
    sys.path.insert(0, local_root)

from syncloud.tools.facade import Facade
from syncloud.app import logger
from syncloud.config.config import PlatformConfig

if __name__ == '__main__':
    www_dir = join(local_root, 'www', '_site')
else:
    www_dir = PlatformConfig().www_root()

logger.init(console=True)

app = Flask(__name__)

@app.route('/server/html/<path:filename>')
def static_file(filename):
    return send_from_directory(www_dir, filename)

@app.route("/server/rest/id", methods=["GET"])
def id():
    return jsonify(success=True, message='', data=convertible.to_dict(Facade().id())), 200

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
