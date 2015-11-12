import traceback
from os import path
from os.path import join
from syncloud_app import logger
from syncloud_platform.sam.stub import SamStub
from syncloud_platform.tools.scripts import run_script


def trigger_app_event_disk(apps_root):
    __trigger_app_event(apps_root, 'on_disk_change.py')


def trigger_app_event_domain(apps_root):
    __trigger_app_event(apps_root, 'on_domain_change.py')


def __trigger_app_event(apps_root, event_script):
    sam = SamStub()
    log = logger.get_logger('events')

    for app in sam.installed_all_apps():
        app_id = app.app.id
        app_event_script = join(apps_root, app_id, 'bin', event_script)
        if path.isfile(app_event_script):
            log.info('executing {0}'.format(app_event_script))
            try:
                run_script(app_event_script)
            except Exception, e:
                log.error('error in script, error: {0}'.format(e.message))
                traceback.print_exc()
        else:
            log.info('{0} not found'.format(app_event_script))
