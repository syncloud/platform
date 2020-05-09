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
            try:
                info = check_output('snap info {0}'.format(app_id), shell=True).decode()
                command = '{0}.{1}'.format(app_id, action)
                if command in info:
                    self.log.info('executing {0}: {1}'.format(app_id, action))
                    output = check_output('snap run {0}'.format(command), shell=True).decode()
                    print(output)
            except CalledProcessError as e:
                self.log.error('event error: {0}'.format(e.output))
                self.log.error(traceback.format_exc())

