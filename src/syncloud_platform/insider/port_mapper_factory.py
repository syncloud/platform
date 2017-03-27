import time

from syncloud_app import logger

from syncloud_platform.insider.config import Port
from syncloud_platform.insider.util import port_to_protocol, is_web_port


class PortMapperFactory:
    def __init__(self, nat_pmp_port_mapper, upnp_port_mapper):
        self.nat_pmp_port_mapper = nat_pmp_port_mapper
        self.upnp_port_mapper = upnp_port_mapper
        self.log = logger.get_logger('port_mapper_factory')

    def check_mapper(self, mapper):
        
        try:
            ip = self.retrying_get_external_address(mapper)
            if ip is None or ip == '':
                raise Exception("Returned bad ip address: {0}".format(ip))
            self.log.warn('{0} mapper is working, returned external ip: {1}'.format(mapper.name(), ip))
            return mapper
        except Exception as e:
            self.log.warn('{0} mapper failed, message: {1}, {2}'.format(mapper.name(), repr(e), vars(e)))
        return None

    def retrying_get_external_address(self, mapper):
        retry = 0
        retries = 5
        ip = mapper.external_ip()
        while not ip and retry < retries:
            retry += 1
            self.log.info('retrying external ip: {0} / {1}'.format(retry, retries))
            time.sleep(1)
            ip = mapper.external_ip()
        return ip

    def provide_mapper(self):
        mapper = self.check_mapper(self.nat_pmp_port_mapper)
        if mapper is not None:
            return mapper
        mapper = self.check_mapper(self.upnp_port_mapper)
        if mapper is not None:
            return mapper
        self.log.error('None of mappers are working')
        return None
