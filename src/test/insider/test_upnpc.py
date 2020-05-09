import pytest

from syncloud_platform.insider.upnpc import UpnpPortMapper, Mapping
from test.insider.inmemory_upnp import InMemoryUPnP


def test_port_free():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')

    mapper = UpnpPortMapper(upnp)
    mapper.add_mapping(80, 80, 'TCP')

    assert len(upnp.mappings) == 1
    assert upnp.mappings[0].external_port == 80
    assert upnp.mappings[0].local_ip == upnp.lanaddr


def test_port_taken():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.mappings = [
        Mapping(80, 'TCP', '3.3.3.3', 80, '', True, '1.1.1.1', '')
    ]

    mapper = UpnpPortMapper(upnp, lower_limit=80)
    mapper.add_mapping(80, 80, 'TCP')

    assert len(upnp.mappings) == 2

    assert upnp.by_external_port(80).local_port == 80
    assert upnp.by_external_port(80).local_ip == '3.3.3.3'

    assert upnp.by_external_port(81).local_port == 80
    assert upnp.by_external_port(81).local_ip == '2.2.2.2'


def test_multiple_external_ports_with_preferred_cleanup():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.mappings = [
        Mapping(80, 'TCP', '2.2.2.2', 80, '', True, '1.1.1.1', ''),
        Mapping(81, 'TCP', '2.2.2.2', 80, '', True, '1.1.1.1', '')
        ]

    mapper = UpnpPortMapper(upnp, lower_limit=80)
    mapper.add_mapping(80, 80, 'TCP')

    assert len(upnp.mappings) == 1

    assert upnp.by_external_port(80).local_port == 80
    assert upnp.by_external_port(80).local_ip == '2.2.2.2'


def test_multiple_external_ports_without_preferred_cleanup_all_except_lowest_port():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.mappings = [
        Mapping(81, 'TCP', '2.2.2.2', 80, '', True, '1.1.1.1', ''),
        Mapping(82, 'TCP', '2.2.2.2', 80, '', True, '1.1.1.1', '')
        ]

    mapper = UpnpPortMapper(upnp, lower_limit=80)
    mapper.add_mapping(80, 80, 'TCP')

    assert len(upnp.mappings) == 1

    assert upnp.by_external_port(81).local_port == 80
    assert upnp.by_external_port(81).local_ip == '2.2.2.2'


def test_fail_to_add():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.fail_on_external_port_with(80, Exception('Failed'))

    mapper = UpnpPortMapper(upnp, lower_limit=81)
    mapper.add_mapping(80, 80, 'TCP')

    assert len(upnp.mappings) == 1

    assert upnp.by_external_port(81).local_port == 80
    assert upnp.by_external_port(81).local_ip == '2.2.2.2'


def test_fail_to_remove():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    upnp.mappings = [Mapping(80, 'TCP', '2.2.2.2', 80, '', True, '1.1.1.1', '')]

    upnp.fail_on_external_port_with(80, Exception('Failed'))

    mapper = UpnpPortMapper(upnp)
    mapper.remove_mapping(80, 80, 'TCP')

    assert len(upnp.mappings) == 1


def test_above_limit_fail_attempts():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    for port in range(1, 6):
        upnp.fail_on_external_port_with(port, Exception('Failed'))

    mapper = UpnpPortMapper(upnp, fail_attempts=5, lower_limit=1)
    with pytest.raises(Exception) as context:
        mapper.add_mapping(1, 1, 'TCP')

    assert 'Unable' in str(context.value)
    assert len(upnp.mappings) == 0


def test_below_limit_fail_attempts():

    upnp = InMemoryUPnP('1.1.1.1', '2.2.2.2')
    for port in range(1, 4):
        upnp.fail_on_external_port_with(port, Exception('Failed'))

    mapper = UpnpPortMapper(upnp, fail_attempts=5, lower_limit=1)
    mapper.add_mapping(1, 1, 'TCP')

    assert len(upnp.mappings) == 1

    assert upnp.by_external_port(4).local_port == 1
    assert upnp.by_external_port(4).local_ip == '2.2.2.2'
