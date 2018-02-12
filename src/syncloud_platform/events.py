import traceback
from os.path import join, isfile, isdir
from syncloud_app import logger
from syncloud_platform.gaplib.scripts import run_script
from syncloud_platform.application.apppaths import AppPaths
from subprocess import CalledProcessError

action_to_old_script = {
    'post-install': 'post-install',
    'pre-remove': 'pre-remove',
    'storage-change': 'on_disk_change.py',
    'access-change': 'on_domain_change.py'
}


def get_script_path(app_paths, action):
    app_dir = app_paths.get_install_dir()
    hooks_path = join(app_dir, 'hooks')
    if isdir(hooks_path):
        return (join(hooks_path, action+'.py'), True)
    else:
        hook_script = action_to_old_script[action]
        return (join(app_dir, 'bin', hook_script), False)


def run_hook_script(app_paths, action):
    app_event_script, add_location_to_sys_path = get_script_path(app_paths, action)
    log = logger.get_logger('events')
    if isfile(app_event_script):
        log.info('executing {0}'.format(app_event_script))
        try:
            run_script(app_event_script, add_location_to_sys_path)
        except CalledProcessError, e:
            log.error('error in script: {0}'.format(e.output))
            log.error(traceback.format_exc())
            if action == 'post-install':
                raise e
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
            app_paths = AppPaths(app_id, self.platform_config)
            run_hook_script(app_paths, action)
