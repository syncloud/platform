from syncloud_platform.insider.upnpc import Mapping


class InMemoryUPnP:
    def __init__(self, externalipaddress, lanaddr):
        self.externalipaddress = externalipaddress
        self.mappings = []
        self.lanaddr = lanaddr

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
        def filter_not(mapping):
            return mapping.external_port == external_port and mapping.protocol == protocol

        self.mappings = filter(filter_not, self.mappings)

    def addportmapping(self, external_port, protocol, local_ip, local_port, description, something):
        self.mappings.append(Mapping(external_port, protocol, local_ip, local_port,
                                     description, True, self.externalipaddress, '1'))
