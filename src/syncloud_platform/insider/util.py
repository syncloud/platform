from syncloudlib.error import PassthroughJsonError
from syncloudlib.json import convertible
from syncloudlib import logger


def check_http_error(response):
    log = logger.get_logger('util')
    if not response.status_code == 200:
        log.error(response.text)
        error = convertible.from_json(response.text)
        raise PassthroughJsonError(error.message, response.text)

