import natpmp
from syncloudlib import logger


class NatPmpPortMapper:

    def __init__(self):
        self.logger = logger.get_logger('NatPmpPortMapper')

    def name(self):
        return 'NAT-PMP'

    def external_ip(self):
        external_ip = natpmp.get_public_address()
        self.logger.info('ip: {0}'.format(external_ip))
        return external_ip

    def add_mapping(self, local_port, external_port, protocol):

        response = natpmp.map_port(protocol_from_string(protocol), external_port, local_port)
        return response.public_port

    def remove_mapping(self, local_port, external_port, protocol):
        natpmp.map_port(protocol_from_string(protocol), external_port, local_port, lifetime=0)


def protocol_from_string(protocol):
    return natpmp.NATPMP_PROTOCOL_TCP if protocol == 'TCP' else natpmp.NATPMP_PROTOCOL_UDP
