import itertools
from subprocess import check_output, CalledProcessError
from miniupnpc import UPnP

from syncloud_app import logger


def port_open_on_router(ip, port):
    try:
        check_output('nc -z -w 1 {0} {1}'.format(ip, port), shell=True)
        return True
    except CalledProcessError, e:
        return False


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
        return self.upnp.externalipaddress()


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
        self.upnp.deleteportmapping(external_port, protocol)

    def add(self, protocol, local_port, external_port, description):
        self.upnp.addportmapping(external_port, protocol, self.upnp.lanaddr, local_port, description, '')


LOWER_LIMIT = 2000
UPPER_LIMIT = 65535
PORTS_TO_TRY = 10


class UpnpPortMapper:

    def __init__(self):
        self.logger = logger.get_logger('UpnpPortMapper')
        self.upnpc = UpnpClient()

    def __find_available_ports(self, existing_ports, local_port, ports_to_try=PORTS_TO_TRY):
        port_range = range(LOWER_LIMIT, UPPER_LIMIT)
        if not local_port in port_range:
            port_range = [local_port] + port_range
        external_ip = self.upnpc.external_ip()
        all_open_ports = (x for x in port_range if not port_open_on_router(external_ip, x) and not x in existing_ports)
        return list(itertools.islice(all_open_ports, 0, ports_to_try))

    def __add_new_mapping(self, local_port):
        existing_ports = self.upnpc.mapped_external_ports('TCP')
        external_ports_to_try = self.__find_available_ports(existing_ports, local_port)
        for external_port in external_ports_to_try:
            try:
                self.logger.debug('mapping {0}->{1} (external->local)'.format(external_port, local_port))
                self.upnpc.add('TCP', local_port, external_port, 'Syncloud')
                return external_port
            except Exception, e:
                self.logger.warn('failed, trying next port: {0}'.format(e.message))
        raise Exception('Unable to add mapping, tried {0} ports'.format(PORTS_TO_TRY))

    def __only_one_mapping(self, external_ports):
        external_ports.sort(reverse=True)
        first_external_port = external_ports.pop()
        for port in external_ports:
            self.upnpc.remove('TCP', port)
        return first_external_port

    def add_mapping(self, local_port):
        external_ports = self.upnpc.get_external_ports('TCP', local_port)
        self.logger.debug("existing router mappings for {0}: {1}".format(local_port, external_ports))
        if len(external_ports) > 0:
            return self.__only_one_mapping(external_ports)
        else:
            return self.__add_new_mapping(local_port)

    def remove_mapping(self, local_port, external_port):
        self.upnpc.remove('TCP', external_port)

    def external_ip(self):
        return self.upnpc.external_ip()
