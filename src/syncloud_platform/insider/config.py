class Port:

    def __init__(self, local_port, external_port):
        self.local_port = local_port
        self.external_port = external_port

    def __str__(self):
        return '{0}->{1}'.format(self.external_port, self.local_port)
\

class Service:

    def __init__(self, name, protocol, type, port, url=None):
        self.name = name
        self.protocol = protocol
        self.type = type
        self.port = port
        self.url = url