import NATPMP
from syncloud_app import logger


class NatPmpPortMapper:

    def __init__(self):
        self.logger = logger.get_logger('NatPmpPortMapper')

    def external_ip(self):
        external_ip = NATPMP.get_public_address()
        self.logger.info('ip: {0}'.format(external_ip))
        return external_ip

    def add_mapping(self, local_port, external_port):
        response = NATPMP.map_port(NATPMP.NATPMP_PROTOCOL_TCP, external_port, local_port)
        return response.public_port

    def remove_mapping(self, local_port, external_port):
        NATPMP.map_port(NATPMP.NATPMP_PROTOCOL_TCP, external_port, local_port, lifetime=0)