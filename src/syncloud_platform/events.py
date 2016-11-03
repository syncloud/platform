import traceback
from os import path
from os.path import join
from syncloud_app import logger
from syncloud_platform.gaplib.scripts import run_script


def run_hook_script(platform_config, hook_script, app_id):
    log = logger.get_logger('events')
    apps_root = platform_config.apps_root()
    app_event_script = join(apps_root, app_id, 'bin', hook_script)
    if path.isfile(app_event_script):
        log.info('executing {0}'.format(app_event_script))
        try:
            run_script(app_event_script)
        except:
            log.error('error in script')
            log.error(traceback.format_exc())
    else:
        log.info('{0} not found'.format(app_event_script))


class EventTrigger:
    def __init__(self, sam, platform_config):
        self.sam = sam
        self.platform_config = platform_config

    def trigger_app_event_disk(self):
        self.__trigger_app_event('on_disk_change.py')

    def trigger_app_event_domain(self):
        self.__trigger_app_event('on_domain_change.py')

    def __trigger_app_event(self, event_script):
        for app in self.sam.installed_all_apps():
            app_id = app.app.id
            run_hook_script(self.platform_config, event_script, app_id)
