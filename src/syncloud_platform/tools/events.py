import traceback
from os import path
from os.path import join
from syncloud_app import logger
from syncloud_platform.gaplib.scripts import run_script


class EventTrigger:
    def __init__(self, sam):
        self.sam = sam

    def trigger_app_event_disk(self, apps_root):
        self.__trigger_app_event(apps_root, 'on_disk_change.py')

    def trigger_app_event_domain(self, apps_root):
        self.__trigger_app_event(apps_root, 'on_domain_change.py')

    def __trigger_app_event(self, apps_root, event_script):
        log = logger.get_logger('events')

        for app in self.sam.installed_all_apps():
            app_id = app.app.id
            app_event_script = join(apps_root, app_id, 'bin', event_script)
            if path.isfile(app_event_script):
                log.info('executing {0}'.format(app_event_script))
                try:
                    run_script(app_event_script)
                except:
                    log.error('error in script')
                    log.error(traceback.format_exc())
            else:
                log.info('{0} not found'.format(app_event_script))
