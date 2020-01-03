import traceback
from syncloudlib import logger
from subprocess import check_output, CalledProcessError


class EventTrigger:
    def __init__(self, installer):
        self.installer = installer
        self.log = logger.get_logger('events')

    def trigger_app_event_disk(self):
        self._trigger_app_event('storage-change')

    def trigger_app_event_domain(self):
        self._trigger_app_event('access-change')

    def _trigger_app_event(self, action):
        for app in self.installer.installed_all_apps():
            app_id = app.app.id
            self.log.info('executing {0}: {1}'.format(app_id, action))
            try:
                output = check_output('snap run {0}.{1}'.format(app_id, action), shell=True)
                print(output)
            except CalledProcessError, e:
                self.log.error('event error: {0}'.format(e.output))
                self.log.error(traceback.format_exc())

