import itertools
from syncloud.app import logger

from syncloud.insider.config import Port


LOWER_LIMIT = 2000
UPPER_LIMIT = 65535
PORTS_TO_TRY = 10


class PortMapper:

    def __init__(self, port_config, upnpc):
        self.port_config = port_config
        self.upnpc = upnpc
        self.logger = logger.get_logger('PortMapper')
        self.external_ip_address = self.upnpc.external_ip()

    def find_available_ports_to_try(self, existing_ports, local_port, ports_to_try=PORTS_TO_TRY):
        port_range = range(LOWER_LIMIT, UPPER_LIMIT)
        if not local_port in port_range:
            port_range = [local_port] + port_range
        all_open_ports = (x for x in port_range if not self._is_external_port_open(x) and not x in existing_ports)
        return list(itertools.islice(all_open_ports, 0, ports_to_try))

    def add(self, local_port):
        existing_ports = self.upnpc.mapped_external_ports("TCP")
        external_ports_to_try = self.find_available_ports_to_try(existing_ports, local_port)
        mapping = None
        for external_port in external_ports_to_try:
            try:
                self.logger.debug("mapping {0}->{1} (external->local)".format(external_port, local_port))
                self.upnpc.add(local_port, external_port)
                mapping = Port(local_port, external_port)
                break
            except Exception, e:
                self.logger.warn('failed, trying next port: {0}'.format(e.message))
        if not mapping:
            raise Exception('Unable to add mapping, tried {0} ports'.format(PORTS_TO_TRY))

        self.port_config.add_or_update(mapping)
        return mapping

    def remove(self, local_port):
        mapping = self.port_config.get(local_port)
        self.upnpc.remove(mapping.external_port)
        self.port_config.remove(local_port)

    def remove_all(self):
        for mapping in self.list():
            self.remove(mapping.local_port)
        self.port_config.remove_all()

    def _is_external_port_open(self, port):
        return self.upnpc.port_open_on_router(self.external_ip_address, port)

    def get(self, local_port):
        return self.port_config.get(local_port)

    def list(self):
        return self.port_config.load()

    def external_ip(self):
        return self.external_ip_address

    def sync(self):
        for mapping in self.list():
            self.sync_one_mapping(mapping)

    def sync_one_mapping(self, mapping):
        external_ports = self.upnpc.get_external_ports("TCP", mapping.local_port)
        self.logger.debug("existing router mappings for {0}: {1}".format(mapping.local_port, external_ports))
        if len(external_ports) > 0:
            external_ports.sort(reverse=True)
            first_external_port = external_ports.pop()
            for port in external_ports:
                self.upnpc.remove(port)
            if first_external_port != mapping.external_port:
                mapping.external_port = first_external_port
                self.port_config.add_or_update(mapping)
        else:
            self.add(mapping.local_port)

    def sync_new_port(self, local_port):
        self.sync_one_mapping(Port(local_port, None))
