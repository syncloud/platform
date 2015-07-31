from syncloud_app import logger

from syncloud_platform.insider.config import Port

from upnpc import UpnpPortMapper
from natpmpc import NatPmpPortMapper

LOWER_LIMIT = 2000
UPPER_LIMIT = 65535
PORTS_TO_TRY = 10


def check_mapper(mapper_name, mapper_type):
    log = logger.get_logger('check_mapper')
    try:
        mapper = mapper_type()
        ip = mapper.external_ip()
        if ip is None or ip == '':
            raise Exception("Returned bad ip address: {0}".format(ip))
        log.warn('{0} mapper is working, returned extrenal ip: {1}'.format(mapper_name, ip))
        return mapper
    except Exception as e:
        log.warn('{0} mapper failed, message: {1}'.format(mapper_name, e.message))
    return None


def provide_mapper():
    log = logger.get_logger('check_mapper')
    mapper = check_mapper('NatPmpPortMapper', NatPmpPortMapper)
    if mapper is not None:
        return mapper
    mapper = check_mapper('UpnpPortMapper', UpnpPortMapper)
    if mapper is not None:
        return mapper
    log.error('None of mappers are working')
    return None


class MockPortMapper:
    def __init__(self, external_ip=None):
        self.__external_ip=external_ip

    def external_ip(self):
        return self.__external_ip

    def add_mapping(self, local_port):
        return local_port

    def remove_mapping(self, local_port, external_port):
        pass


class PortDrill:
    def __init__(self, port_config, port_mapper_provider):
        self.logger = logger.get_logger('PortMapper')
        self.port_config = port_config
        self.port_mapper = port_mapper_provider()

    def remove_all(self):
        for mapping in self.list():
            self.remove(mapping.local_port)
        self.port_config.remove_all()

    def get(self, local_port):
        return self.port_config.get(local_port)

    def list(self):
        return self.port_config.load()

    def external_ip(self):
        return self.port_mapper.external_ip()

    def remove(self, local_port):
        mapping = self.port_config.get(local_port)
        self.port_mapper.remove_mapping(mapping.local_port, mapping.external_port)
        self.port_config.remove(local_port)

    def sync_one_mapping(self, local_port):
        external_port = self.port_mapper.add_mapping(local_port)
        mapping = Port(local_port, external_port)
        self.port_config.add_or_update(mapping)

    def sync_new_port(self, local_port):
        self.sync_one_mapping(local_port)

    def sync(self):
        for mapping in self.list():
            self.sync_one_mapping(mapping.local_port)