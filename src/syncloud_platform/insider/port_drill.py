import time

from syncloud_app import logger

from syncloud_platform.insider.config import Port
from syncloud_platform.insider.port_prober import PortProber
from syncloud_platform.insider.util import port_to_protocol, is_web_port


def check_mapper(mapper):
    log = logger.get_logger('check_mapper')
    try:
        ip = retrial_get_external_address(mapper)
        if ip is None or ip == '':
            raise Exception("Returned bad ip address: {0}".format(ip))
        log.warn('{0} mapper is working, returned external ip: {1}'.format(mapper.name(), ip))
        return mapper
    except Exception as e:
        log.warn('{0} mapper failed, message: {1}, {2}'.format(mapper.name(), repr(e), vars(e)))
    return None


def retrial_get_external_address(mapper):
    log = logger.get_logger('get_external_address')
    retry = 0
    retries = 5
    ip = mapper.external_ip()
    while not ip and retry < retries:
        retry += 1
        log.info('retrying external ip: {0} / {1}'.format(retry, retries))
        time.sleep(1)
        ip = mapper.external_ip()
    return ip


def provide_mapper(nat_pmp_port_mapper, upnp_port_mapper):
    log = logger.get_logger('check_mapper')
    mapper = check_mapper(nat_pmp_port_mapper)
    if mapper is not None:
        return mapper
    mapper = check_mapper(upnp_port_mapper)
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
        if mapping:
            self.port_mapper.remove_mapping(mapping.local_port, mapping.external_port, 'TCP')
            self.port_config.remove(local_port)

    def sync_one_mapping(self, local_port, protocol):

        self.logger.info('Sync one mapping: {0}'.format(local_port))
        port_to_try = local_port
        lower_limit = 10000
        found_external_port = None
        retries = 10
        for i in range(1, retries):
            self.logger.info('Trying {0}'.format(port_to_try))
            external_port = self.port_mapper.add_mapping(local_port, port_to_try, protocol)
            if not is_web_port(local_port):
                found_external_port = external_port
                break
            if self.port_prober.probe_port(external_port, port_to_protocol(local_port)):
                found_external_port = external_port
                break
            self.port_mapper.remove_mapping(local_port, external_port, protocol)

            if port_to_try == local_port:
                port_to_try = lower_limit
            else:
                port_to_try = external_port + 1

        if not found_external_port:
            raise Exception('Unable to add mapping, tried {0} times'.format(retries))

        mapping = Port(local_port, found_external_port)
        self.port_config.add_or_update(mapping)

    def sync_new_port(self, local_port):
        self.sync_one_mapping(local_port, 'TCP')

    def sync(self):
        for mapping in self.list():
            self.sync_one_mapping(mapping.local_port, 'TCP')

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


class PortDrillFactory:
    def __init__(self, user_platform_config, port_config, nat_pmp_port_mapper, upnp_port_mapper):
        self.port_config = port_config
        self.user_platform_config = user_platform_config
        self.nat_pmp_port_mapper = nat_pmp_port_mapper
        self.upnp_port_mapper = upnp_port_mapper

    def get_drill(self, external_access):
        if not external_access:
            return NonePortDrill()
        drill = None
        mapper = provide_mapper(self.nat_pmp_port_mapper, self.upnp_port_mapper)
        if mapper:
            prober = PortProber(
                self.user_platform_config.get_redirect_api_url(),
                self.user_platform_config.get_domain_update_token())
            drill = PortDrill(self.port_config, mapper, prober)
        return drill
