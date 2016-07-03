import itertools
from subprocess import check_output, CalledProcessError
from miniupnpc import UPnP

from syncloud_app import logger
import time


class Mapping:
    def __init__(self, external_port, protocol, local_ip, local_port, description, enabled, remote_ip, lease_time):
        self.external_port = external_port
        self.protocol = protocol
        self.local_ip = local_ip
        self.local_port = local_port
        self.description = description
        self.enabled = enabled
        self.remote_ip = remote_ip
        self.lease_time = lease_time


def to_mapping(m):
    external_port, protocol, local_address, description, enabled_str, remote_ip_str, lease_time = m
    local_ip_str, local_port = local_address
    local_ip = local_ip_str
    if local_ip_str == '': local_ip = None
    remote_ip = remote_ip_str
    if remote_ip_str == '': remote_ip = None
    enabled = False
    if enabled_str == '1': enabled = True
    return Mapping(external_port, protocol, local_ip, local_port, description, enabled, remote_ip, lease_time)


class UpnpClient:
    def __init__(self, upnp):
        self.logger = logger.get_logger('UpnpClient')
        self.upnp = upnp
        self.initialized = False

    def init(self):
        if self.initialized:
            return
        self.logger.info('initializing upnp')
        self.upnp.discover()
        self.info(upnp.devlist)
        self.upnp.selectigd()
        self.initialized = True

    def __run(self, cmd):
        return check_output(cmd, shell=True)

    def external_ip(self):
        external_ip = self.upnp.externalipaddress()
        self.logger.info('ip: {0}'.format(external_ip))
        return external_ip

    def __list(self):
        result = []
        i = 0
        while True:
            p = self.upnp.getgenericportmapping(i)
            if p is None:
                break
            result.append(p)
            i += 1
        return [to_mapping(m) for m in result]

    def mapped_external_ports(self, protocol):
        mappings = self.__list()
        ports = [m.external_port for m in mappings if m.protocol == protocol]
        return ports

    def get_external_ports(self, protocol, local_port):
        mappings = self.__list()
        local_ip = self.upnp.lanaddr
        ports = [m.external_port for m in mappings if m.protocol == protocol and m.local_ip == local_ip and m.local_port == local_port]
        return ports

    def remove(self, protocol, external_port):
        self.logger.info('removing {0} port mapping'.format(external_port))
        self.upnp.deleteportmapping(external_port, protocol)

    def add(self, protocol, local_port, external_port, description):
        self.logger.debug('adding {0} -> {1} port mapping'.format(external_port, local_port))
        self.upnp.addportmapping(external_port, protocol, self.upnp.lanaddr, local_port, description, '')


LOWER_LIMIT = 10000
UPPER_LIMIT = 65535


class UpnpPortMapper:

    def __init__(self, upnp, fail_attempts=50):
        self.fail_attempts = fail_attempts
        self.logger = logger.get_logger('UpnpPortMapper')
        self.upnp_client = UpnpClient(upnp)

    def name(self):
        return 'UpnpPortMapper'

    def upnpc(self):
        self.upnp_client.init()
        return self.upnp_client

    def __find_available_ports(self, existing_ports, external_port):
        port_range = range(external_port, UPPER_LIMIT)
        available_ports = [x for x in port_range if x not in existing_ports]
        return available_ports[0:self.fail_attempts]

    def __add_new_mapping(self, local_port, external_port, protocol):
        existing_ports = self.upnpc().mapped_external_ports(protocol)
        external_ports_to_try = self.__find_available_ports(existing_ports, external_port)
        for external_port_to_try in external_ports_to_try:
            try:
                self.logger.info('mapping {0}->{1} (external->local)'.format(external_port_to_try, local_port))
                self.upnpc().add(protocol, local_port, external_port_to_try, 'Syncloud')

                existing_ports = self.upnpc().mapped_external_ports(protocol)
                self.logger.info('ports after mapping {0}'.format(existing_ports))

                return external_port_to_try
            except Exception, e:
                self.logger.error('failed: {0}, {1}'.format(repr(e), vars(e)))

        raise Exception('Unable to add mapping')

    def __only_one_mapping(self, external_ports, protocol):
        external_ports.sort(reverse=True)
        first_external_port = external_ports.pop()
        for port in external_ports:
            self.upnpc().remove(protocol, port)
        return first_external_port

    def add_mapping(self, local_port, external_port, protocol):
        external_ports = self.upnpc().get_external_ports(protocol, local_port)
        self.logger.info("existing router mappings for {0}: {1}".format(local_port, external_ports))
        if len(external_ports) > 0:
            return self.__only_one_mapping(external_ports, protocol)
        else:
            return self.__add_new_mapping(local_port, external_port, protocol)

    def remove_mapping(self, local_port, external_port, protocol):
        try:
            self.upnpc().remove(protocol, external_port)
        except Exception, e:
            self.logger.warn('unable to remove port {0}, probably does not exist anymore, error: {1}, {2}'.format(
                             external_port, repr(e), vars(e)))

    def external_ip(self):

        retry = 0
        retries = 5
        ip = self.upnpc().external_ip()
        while not ip and retry < retries:
            retry += 1
            self.logger.info('retrying external ip: {0} / {1}'.format(retry, retries))
            time.sleep(1)
            ip = self.upnpc().external_ip()
        return ip
