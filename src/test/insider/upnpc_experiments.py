from miniupnpc import UPnP

from syncloud_platform.insider.upnpc import UpnpPortMapper, UpnpClient


def test_many_ports_until_fail():
    mapper = UpnpPortMapper(UpnpClient(UPnP()))
    base_port = 11000
    for port in range(base_port, base_port + 50, 1):
        mapper.add_mapping(port, port, 'TCP')
