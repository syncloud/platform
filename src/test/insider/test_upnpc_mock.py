import logging
from syncloud.app import logger
from syncloud.insider import upnpc_mock

logger.init(level=logging.DEBUG, console=True)


def test_mock():
    upnpc = upnpc_mock.Upnpc()
    assert len(upnpc.mapped_external_ports('TCP')) == 0
    upnpc.add(80, 10000)
    assert len(upnpc.mapped_external_ports('TCP')) == 1
    assert len(upnpc.get_external_ports('TCP', 80)) == 1
    upnpc.add(81, 10001)
    assert len(upnpc.mapped_external_ports('TCP')) == 2
    assert len(upnpc.get_external_ports('TCP', 80)) == 1
    assert len(upnpc.get_external_ports('TCP', 81)) == 1
    upnpc.remove(10000)
    assert len(upnpc.mapped_external_ports('TCP')) == 1
    assert len(upnpc.get_external_ports('TCP', 80)) == 0
    assert len(upnpc.get_external_ports('TCP', 81)) == 1
    upnpc.remove(10001)
    assert len(upnpc.mapped_external_ports('TCP')) == 0
    assert len(upnpc.get_external_ports('TCP', 80)) == 0
    assert len(upnpc.get_external_ports('TCP', 81)) == 0
