from functools import update_wrapper
from flask import make_response, redirect, request
from syncloud_platform.config.config import PlatformConfig, PlatformUserConfig

def nocache(f):
    def new_func(*args, **kwargs):
        resp = make_response(f(*args, **kwargs))
        # resp.cache_control.no_cache = True
        resp.headers['Cache-Control'] = 'no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0'
        return resp
    return update_wrapper(new_func, f)

def redirect_if_not_activated(f):
    def new_func(*args, **kwargs):
        resp = make_response(f(*args, **kwargs))
        platform_config = PlatformConfig()
        if not PlatformUserConfig(platform_config.get_user_config()).is_activated():
            return redirect('{0}://{1}:81'.format(request.scheme, request.host))
        else:
            return resp
    return update_wrapper(new_func, f)
