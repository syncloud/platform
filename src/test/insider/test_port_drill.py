from syncloud_platform.insider.config import Port
from syncloud_platform.insider.port_drill import PortDrill
from test.insider.helpers import get_port_config


def test_port_drill():
    port_config = get_port_config([Port(80, 80, 'TCP')])
    port_drill = PortDrill(port_config, MockPortMapper(external_ip='192.167.44.52'), MockPortProber())
    port_drill.sync_existing_ports()
    mapping = port_drill.get(80, 'TCP')
      
    assert mapping.external_port == 80
    assert port_drill.external_ip() == '192.167.44.52'
    

class MockPortMapper:
    def __init__(self, external_ip=None):
        self.__external_ip = external_ip

    def external_ip(self):
        return self.__external_ip

    def add_mapping(self, local_port, external_port, protocol):
        return external_port

    def remove_mapping(self, local_port, external_port, protocol):
        pass


class MockPortProber:

    def probe_port(self, port, protocol):
        return True
