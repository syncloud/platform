from syncloud_app import logger


class Upnpc():

    def __init__(self):
        self.logger = logger.get_logger('mock.Upnpc')
        self.logger.info('initializing mock')
        self.mappings = dict()

    def external_ip(self):
        return None

    def mapped_external_ports(self, protocol):
        self.logger.debug('mapped_external_ports({0})'.format(protocol))
        return [external for (external, local) in self.mappings.iteritems()]

    def get_external_ports(self, protocol, local_port):
        self.logger.debug('get_external_ports({0}, {1})'.format(protocol, local_port))
        return [external for (external, local) in self.mappings.iteritems() if local == local_port]

    def remove(self, external_port):
        self.logger.debug('remove({0})'.format(external_port))
        del self.mappings[external_port]

    def add(self, local_port, external_port):
        self.logger.debug('add({0}, {1})'.format(local_port, external_port))
        self.mappings[external_port] = local_port

    def port_open_on_router(self, ip, port):
        return False