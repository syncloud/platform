import logging

from syncloud_app import logger

from syncloud_platform.insider.upnpc import UpnpPortMapper, UpnpClient, Mapping
from test.insider.inmemory_upnp import InMemoryUPnP

logger.init(level=logging.DEBUG, console=True)


def test_port_free():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')

    mapper = UpnpPortMapper(UpnpClient(upnp))
    mapper.add_mapping(80, 80)

    assert len(upnp.mappings) == 1
    assert upnp.mappings[0].external_port == 80
    assert upnp.mappings[0].local_ip == upnp.lanaddr


def test_port_taken():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.mappings = [Mapping(80, 'TCP', '3.3.3.3', 80, '', True, '1.1.1.1', '')]

    mapper = UpnpPortMapper(UpnpClient(upnp))
    mapper.add_mapping(80, 80)

    assert len(upnp.mappings) == 2

    assert upnp.by_external_port(80).local_port == 80
    assert upnp.by_external_port(80).local_ip == '3.3.3.3'

    assert upnp.by_external_port(81).local_port == 80
    assert upnp.by_external_port(81).local_ip == '2.2.2.2'


def test_fail_to_add():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.fail_on_external_port_with(80, Exception('Failed'))

    mapper = UpnpPortMapper(UpnpClient(upnp))
    mapper.add_mapping(80, 80)

    assert len(upnp.mappings) == 1

    assert upnp.by_external_port(81).local_port == 80
    assert upnp.by_external_port(81).local_ip == '2.2.2.2'


def test_fail_to_remove():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.mappings = [Mapping(80, 'TCP', '2.2.2.2', 80, '', True, '1.1.1.1', '')]

    upnp.fail_on_external_port_with(80, Exception('Failed'))

    mapper = UpnpPortMapper(UpnpClient(upnp))
    mapper.remove_mapping(80, 80)

    assert len(upnp.mappings) == 1
