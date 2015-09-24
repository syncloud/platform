from urlparse import urljoin
import requests
from syncloud_app import logger

from syncloud_platform.insider.config import Port

from upnpc import UpnpPortMapper
from natpmpc import NatPmpPortMapper


def check_mapper(mapper_type):
    log = logger.get_logger('check_mapper')
    mapper_name = mapper_type.__name__
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
    mapper = check_mapper(NatPmpPortMapper)
    if mapper is not None:
        return mapper
    mapper = check_mapper(UpnpPortMapper)
    if mapper is not None:
        return mapper
    log.error('None of mappers are working')
    return None


class PortDrill:
    def __init__(self, port_config, port_mapper, port_prober):
        self.port_prober = port_prober
        self.logger = logger.get_logger('PortDrill')
        self.port_config = port_config
        self.port_mapper = port_mapper

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
        self.logger.info('Sync one mapping: {0}'.format(local_port))
        port_to_try = local_port
        lower_limit = 10000
        found_external_port = None
        for i in range(1, 10):
            self.logger.info('Trying {0}'.format(port_to_try))
            external_port = self.port_mapper.add_mapping(local_port, port_to_try)
            if self.port_prober.probe_port(external_port):
                found_external_port = external_port
                break
            else:
                self.port_mapper.remove_mapping(local_port, external_port)

            if port_to_try == local_port:
                port_to_try = lower_limit
            else:
                port_to_try = external_port + 1

        if not found_external_port:
            raise Exception('Unable to add mapping')

        mapping = Port(local_port, found_external_port)
        self.port_config.add_or_update(mapping)

    def sync_new_port(self, local_port):
        self.sync_one_mapping(local_port)

    def sync(self):
        for mapping in self.list():
            self.sync_one_mapping(mapping.local_port)

    def available(self):
        return self.port_mapper is not None



class NonePortDrill:
    def __init__(self):
        self.logger = logger.get_logger('NonePortDrill')

    def remove_all(self):
        pass

    def get(self, local_port):
        return Port(local_port, None)

    def list(self):
        return []

    def external_ip(self):
        return None

    def remove(self, local_port):
        pass

    def sync_one_mapping(self, local_port):
        pass

    def sync_new_port(self, local_port):
        pass

    def sync(self):
        pass

    def available(self):
        return False
