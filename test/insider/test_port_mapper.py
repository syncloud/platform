import logging
import unittest

from mock import MagicMock

from syncloud.insider.port_mapper import PortMapper
from syncloud.insider.config import Port
from test.insider.helpers import get_port_config


logging.basicConfig(level=logging.DEBUG)


class FailingUpnpc():
    def __init__(self, to_fail):
        self.to_fail = to_fail
        self.failed = 0

    def add(self, a, b):
        if self.failed < self.to_fail:
            self.failed += 1
            raise Exception('fail')
        else:
            return 'good'


def test_add_success():
    port_config = get_port_config([])

    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)

    mapper = PortMapper(port_config, upnpc)
    mapping = mapper.add(80)

    assert mapping is not None
    assert 80 == mapping.local_port
    assert 80 == mapping.external_port
    assert upnpc.add.called

    read = port_config.get(80)
    assert 80 == read.external_port


def test_add_success_after_first_failed():
    port_config = get_port_config([])

    upnpc = FailingUpnpc(2)
    upnpc.port_open_on_router = MagicMock(return_value=False)
    upnpc.external_ip = MagicMock()
    upnpc.mapped_external_ports = MagicMock()

    mapper = PortMapper(port_config, upnpc)
    mapping = mapper.add(80)

    assert mapping is not None
    assert 80 == mapping.local_port
    assert 2001 == mapping.external_port

    read = port_config.get(80)
    assert 2001 == read.external_port


def test_sync_adds_port():
    port_config = get_port_config([Port(80, 10001)])
    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)

    mapper = PortMapper(port_config, upnpc)
    mapper.sync()

    read = port_config.get(80)
    assert 80 == read.external_port

    assert upnpc.add.called


def test_sync_cleans_duplicate_mappings():
    port_config = get_port_config([Port(80, 2000)])

    upnpc = MagicMock()
    upnpc.get_external_ports = MagicMock(return_value=[2000, 2001])

    mapper = PortMapper(port_config, upnpc)
    mapper.sync()

    read = port_config.get(80)
    assert 2000 == read.external_port

    upnpc.remove.assert_called_with(2001)


def test_sync_new_ports():
    port_config = get_port_config([])

    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)

    mapper = PortMapper(port_config, upnpc)
    mapper.sync_new_port(80)

    read = port_config.get(80)
    assert 80 == read.external_port

    upnpc.add.assert_called_with(80, 80)


def test_first_gap():
    port_config = get_port_config([])
    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)
    mapper = PortMapper(port_config, upnpc)
    existing_ports = [2000, 2001, 2003]
    ports_to_try = mapper.find_available_ports_to_try(existing_ports, 2000, 3)
    assert [2002, 2004, 2005] == ports_to_try


def test_no_gap():
    port_config = get_port_config([])
    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)
    mapper = PortMapper(port_config, upnpc)
    existing_ports = [2000, 2001, 2002]
    assert [2003, 2004, 2005] == mapper.find_available_ports_to_try(existing_ports, 2000, 3)


def test_no_existing():
    port_config = get_port_config([])
    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)
    mapper = PortMapper(port_config, upnpc)
    existing_ports = []
    assert [2000, 2001, 2002] == mapper.find_available_ports_to_try(existing_ports, 2000, 3)


def test_open_port():
    port_config = get_port_config([])
    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(side_effect=[True, True, False, False, False])
    mapper = PortMapper(port_config, upnpc)
    existing_ports = []
    assert [2001, 2002, 2003] == mapper.find_available_ports_to_try(existing_ports, 1999, 3)


def test_local_port():
    port_config = get_port_config([])
    upnpc = MagicMock()
    upnpc.port_open_on_router = MagicMock(return_value=False)
    mapper = PortMapper(port_config, upnpc)
    existing_ports = [2000, 2001, 2003]
    assert [80, 2002, 2004] == mapper.find_available_ports_to_try(existing_ports, 80, 3)


def test_lazy_external_ip_expensive_call():
    port_config = get_port_config([])

    upnpc = MagicMock()
    upnpc.external_ip = MagicMock()
    mapper = PortMapper(port_config, upnpc)

    assert not upnpc.external_ip.called
    
    mapper.external_ip()
    mapper.external_ip()

    assert upnpc.external_ip.call_count == 1
