import sys
import convertible
import inspect
import logger
import logging
import traceback


def call(func, kwargs):
    argspec = inspect.getargspec(func)
    params = kwargs
    if not argspec.keywords:
        params = dict()
        for key in kwargs.keys():
            if key in argspec.args:
                params[key] = kwargs.get(key)
    return func(**params)


def respond(result, message=None, success=True):
    convertible.pretty_print()
    response = dict(success=success, message=message, data=result)
    response_json = convertible.to_json(response)
    convertible.ugly_print()
    print(response_json)


def run(runner, debug=False, text=False, exit_code_on_error=1):
    log = logger.get_logger('app')
    try:
        data = runner()
        if text:
            print(data)
        else:
            respond(data)
    except Exception, e:
        log.exception(e)
        respond(e, str(e), success=False)
        exit(exit_code_on_error)


def execute(obj, args, exit_code_on_error=0):
    method = getattr(obj, args.action)
    text = False
    if hasattr(args, 'text'):
        text = args.text
    run(lambda: call(method, vars(args)), args.debug, text, exit_code_on_error)


# Leaving for backward compatibility until all apps upgraded
# Use directly logger.init()
def init_log(name, args):
    console = True if args.debug else False
    level = logging.DEBUG if args.debug else logging.INFO
    logger.init(level, console, '/tmp/{}.log'.format(name))
