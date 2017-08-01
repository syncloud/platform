from syncloud_app import logger

from syncloud_platform.insider.util import is_web_port


class ManualPortMapper:

    def __init__(self, manual_public_ip, manual_public_web_port):
        self.manual_public_web_port = manual_public_web_port
        self.manual_public_ip = manual_public_ip
        self.logger = logger.get_logger('ManualPortMapper')

    def name(self):
        return 'Manual'

    def external_ip(self):
        return self.manual_public_ip

    def add_mapping(self, local_port, external_port, protocol):
        self.logger.warn('adding port mapping is not available in manual mode {0}:{1}'.format(
            self.manual_public_ip, self.manual_public_web_port))
        if is_web_port(local_port):
            return self.manual_public_web_port
        else:
            return external_port

    def remove_mapping(self, local_port, external_port, protocol):
        self.logger.warn('removing port mapping is not available in manual mode')
