import itertools
from subprocess import check_output, CalledProcessError
from miniupnpc import UPnP

from syncloud_app import logger

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
    def __init__(self):
        self.logger = logger.get_logger('UpnpClient')
        self.upnp = UPnP()
        self.upnp.discover()
        self.upnp.selectigd()

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
        local_ip = self.upnp.lanaddr
        ports = [m.external_port for m in mappings if m.protocol == protocol and m.local_ip == local_ip]
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

    def __init__(self):
        self.logger = logger.get_logger('UpnpPortMapper')
        self.__upnpc = None

    def upnpc(self):
        if self.__upnpc is None:
            self.__upnpc = UpnpClient()
        return self.__upnpc

    def __find_available_ports(self, existing_ports, external_port):
        port_range = range(external_port, UPPER_LIMIT)
        return (x for x in port_range if x not in existing_ports)

    def __add_new_mapping(self, local_port, external_port):
        existing_ports = self.upnpc().mapped_external_ports('TCP')
        external_ports_to_try = self.__find_available_ports(existing_ports, external_port)
        for external_port_to_try in external_ports_to_try:
            try:
                self.logger.debug('mapping {0}->{1} (external->local)'.format(external_port_to_try, local_port))
                self.upnpc().add('TCP', local_port, external_port_to_try, 'Syncloud')
                return external_port_to_try
            except Exception, e:
                self.logger.debug('failed, trying next port: {0}'.format(e.message))

        raise Exception('Unable to add mapping')

    def __only_one_mapping(self, external_ports):
        external_ports.sort(reverse=True)
        first_external_port = external_ports.pop()
        for port in external_ports:
            self.upnpc().remove('TCP', port)
        return first_external_port

    def add_mapping(self, local_port, external_port):
        external_ports = self.upnpc().get_external_ports('TCP', local_port)
        self.logger.info("existing router mappings for {0}: {1}".format(local_port, external_ports))
        if len(external_ports) > 0:
            return self.__only_one_mapping(external_ports)
        else:
            return self.__add_new_mapping(local_port, external_port)

    def remove_mapping(self, local_port, external_port):
        self.upnpc().remove('TCP', external_port)

    def external_ip(self):
        return self.upnpc().external_ip()
