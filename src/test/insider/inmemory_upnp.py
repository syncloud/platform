from syncloud_platform.insider.upnpc import Mapping


class InMemoryUPnP:
    def __init__(self, externalipaddress, lanaddr):
        self.externalipaddress = externalipaddress
        self.mappings = []
        self.lanaddr = lanaddr
        self.external_port_to_fail = dict()
        self.devlist = []

    def fail_on_external_port_with(self, external_port, e):
        self.external_port_to_fail[external_port] = e

    def discover(self):
        pass

    def selectigd(self):
        pass

    def by_external_port(self, external_port):
        for mapping in self.mappings:
            if mapping.external_port == external_port:
                return mapping

    def externalipaddress(self):
        return self.externalipaddress

    def getgenericportmapping(self, index):
        if index < len(self.mappings):
            mapping = self.mappings[index]
            return mapping.external_port, mapping.protocol, (mapping.local_ip, mapping.local_port), \
                   mapping.description, mapping.enabled, mapping.remote_ip, mapping.lease_time
        return None

    def deleteportmapping(self, external_port, protocol):

        if external_port in self.external_port_to_fail:
            raise self.external_port_to_fail[external_port]

        def filter_not(mapping):
            return mapping.external_port == external_port and mapping.protocol == protocol

        self.mappings = filter(filter_not, self.mappings)

    def addportmapping(self, external_port, protocol, local_ip, local_port, description, something):

        if external_port in self.external_port_to_fail:
            raise self.external_port_to_fail[external_port]

        self.mappings.append(Mapping(external_port, protocol, local_ip, local_port,
                                     description, True, self.externalipaddress, '1'))
