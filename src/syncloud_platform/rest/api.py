import sys
import traceback

from flask import Flask, jsonify, request, Response
from syncloudlib.error import PassthroughJsonError

app = Flask(__name__)



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
