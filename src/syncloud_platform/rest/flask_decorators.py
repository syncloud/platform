from functools import update_wrapper
from flask import make_response, redirect
from syncloud_platform.injector import get_injector
from syncloudlib.logger import get_logger


def nocache(f):
    def new_func(*args, **kwargs):
        resp = make_response(f(*args, **kwargs))
        # resp.cache_control.no_cache = True
        resp.headers['Cache-Control'] = 'no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0'
        return resp
    return update_wrapper(new_func, f)


def redirect_if_not_activated(f):

    def new_func(*args, **kwargs):
        if _is_activated():
            return make_response(f(*args, **kwargs))
        else:
            return redirect('/activate.html')

    return update_wrapper(new_func, f)


def redirect_if_activated(f):
    def new_func(*args, **kwargs):
        if _is_activated():
            return redirect('/')
        else:
            return make_response(f(*args, **kwargs))

    return update_wrapper(new_func, f)


def _is_activated():
    log = get_logger('activated_check')
    platform_user_config = get_injector().user_platform_config
    activated = False
    try:
        activated = platform_user_config.is_activated()
    except Exception as e:
        log.error('unable to verify activation status, assume it is not activated, {0}'.format(str(e)))

    return activated