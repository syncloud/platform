from syncloud_app.main import PassthroughJsonError
import convertible
from syncloud_app import logger


def check_http_error(response):
    log = logger.get_logger('util')
    if not response.status_code == 200:
        log.error(response.text)
        error = convertible.from_json(response.text)
        raise PassthroughJsonError(error.message, response.text)

def port_to_protocol(port):
    if port == 443:
        return 'https'
    return 'http'

def protocol_to_port(protocol):
    if protocol == 'https':
        return 443
    return 80
