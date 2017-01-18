import traceback
from os.path import join, isfile, isdir
from syncloud_app import logger
from syncloud_platform.gaplib.scripts import run_script


action_to_old_script = {
    'post-install': 'post-install',
    'pre-remove': 'pre-remove',
    'storage-change': 'on_disk_change.py',
    'access-change': 'on_domain_change.py'
}


def get_script_path(apps_root, app_id, action):
    hooks_path = join(apps_root, app_id, 'hooks')
    if isdir(hooks_path):
        return (join(hooks_path, action+'.py'), True)
    else:
        hook_script = action_to_old_script[action]
        return (join(apps_root, app_id, 'bin', hook_script), False)


def run_hook_script(apps_root, app_id, action):
    app_event_script, add_location_to_sys_path = get_script_path(apps_root, app_id, action)
    log = logger.get_logger('events')
    if isfile(app_event_script):
        log.info('executing {0}'.format(app_event_script))
        try:
            run_script(app_event_script, add_location_to_sys_path)
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
        self.__trigger_app_event('storage-change')

    def trigger_app_event_domain(self):
        self.__trigger_app_event('access-change')

    def __trigger_app_event(self, action):
        for app in self.sam.installed_all_apps():
            app_id = app.app.id
            run_hook_script(self.platform_config.get_hooks_root(), app_id, action)
