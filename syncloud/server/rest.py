import traceback
import convertible
from flask import Flask, jsonify, send_from_directory
from os.path import dirname, join
import sys
if __name__ == '__main__':
    sys.path.insert(0, join(dirname(__file__), '..', '..'))
from syncloud.insider.dns import Endpoint
from syncloud.server.model import Site
from syncloud.app import logger
from syncloud.insider.config import Service, os
from syncloud.insider.facade import get_insider

logger.init(console=True)

app = Flask(__name__)
www_dir = '/var/www/syncloud-server'
mock_sites = None
html_prefix = ''
rest_prefix = ''

if __name__ == '__main__':
    html_prefix = '/server/html'
    rest_prefix = '/server/rest'

@app.route(html_prefix + '/<path:filename>')
def static_file(filename):
    return send_from_directory(www_dir, filename)


@app.route('/')
def index():
    return static_file('index.html')


@app.route(rest_prefix + "/sites", methods=["GET"])
def site_list():
    if mock_sites:
        endpoints = mock_sites
    else:
        endpoints = filter_websites(get_insider().endpoints())

    return jsonify(sites=convertible.to_dict(map(Site, endpoints))), 200


@app.errorhandler(Exception)
def handle_exception(error):
    print '-'*60
    traceback.print_exc(file=sys.stdout)
    print '-'*60
    response = jsonify(message=error.message)
    status_code = 500
    return response, status_code


def filter_websites(endpoints):
    return [endpoint for endpoint in endpoints if endpoint.service.protocol in ["http", "https"] and endpoint.service.name != "server"]


if __name__ == '__main__':
    www_dir = join(dirname(__file__), '..', '..', 'www')
    mock_sites = [
        Endpoint(Service("image-ci", "http", "type", "80", "image-ci"), 'localhost', 8181),
        Endpoint(Service("owncloud", "https", "type", "443", "owncloud"), 'localhost', 8282)]
    app.run(host='0.0.0.0', debug=True, port=5001)