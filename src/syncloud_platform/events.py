from syncloudlib import logger

from syncloud_platform.rest.backend_proxy import backend_request


class EventTrigger:
    def __init__(self, installer):
        self.installer = installer
        self.log = logger.get_logger('events')

    def trigger_app_event_disk(self):
        self._trigger_app_event('storage-change')

    def trigger_app_event_domain(self):
        self._trigger_app_event('access-change')

    def _trigger_app_event(self, action):
        try:
            self.log.info('event trigger: {0}'.format(action))
            backend_request("POST", "/event/trigger", {"event": action})
        except Exception as e:
            self.log.error('event error: {0}'.format(str(e)))
