from syncloud_app import logger

from syncloud.insider.config import Port
from syncloud.insider.upnpc import UpnpPortMapper

LOWER_LIMIT = 2000
UPPER_LIMIT = 65535
PORTS_TO_TRY = 10


class PortMapper:

    def __init__(self, port_config, cmd):
        self.logger = logger.get_logger('PortMapper')
        self.port_config = port_config
        self.port_mapper = UpnpPortMapper(cmd)

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
        self.sync_one_mapping(Port(local_port, None))

    def sync(self):
        for mapping in self.list():
            self.sync_one_mapping(mapping)
