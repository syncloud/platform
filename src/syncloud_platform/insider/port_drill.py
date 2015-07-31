from syncloud_app import logger

from syncloud_platform.insider.config import Port

from upnpc import UpnpPortMapper
from natpmpc import NatPmpPortMapper

LOWER_LIMIT = 2000
UPPER_LIMIT = 65535
PORTS_TO_TRY = 10


class MockPortMapper:
    def __init__(self, external_ip=None):
        self.__external_ip=external_ip

    def external_ip(self):
        return self.__external_ip

    def add_mapping(self, local_port):
        return local_port

    def remove_mapping(self, local_port, external_port):
        pass


class PortMapperProvider:
    def __init__(self):
        self.logger = logger.get_logger('PortMapperProvider')

    def __call__(self):
        mapper = self.__check_mapper('NAT-PMP', NatPmpPortMapper)
        if mapper is not None:
            return mapper
        mapper = self.__check_mapper('UPnP', UpnpPortMapper)
        if mapper is not None:
            return mapper
        self.logger.error('None of mappers are working')

    def __check_mapper(self, mapper_name, mapper_type):
        try:
            mapper = mapper_type()
            ip = mapper.external_ip()
            if ip is None:
                raise Exception("No external ip address returned")
            self.logger.warn('{0} mapper is working, returned extrenal ip: {1}'.format(mapper_name, ip))
            return mapper
        except Exception as e:
            self.logger.warn('{0} mapper failed, message: {1}'.format(mapper_name, e.message))
        return None



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