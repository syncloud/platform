import logging

import pytest

from syncloud_app import logger

from syncloud_platform.insider.upnpc import UpnpPortMapper
from syncloud_platform.insider.natpmpc import NatPmpPortMapper
from test.insider.http import SomeHttpServer, wait_http, wait_http_cant_connect

logger.init(level=logging.DEBUG, console=True)

@pytest.fixture(scope="module")
def http_server(request):
    server = SomeHttpServer(8080)
    server.start()

    def fin():
        server.stop()
    request.addfinalizer(fin)
    return server

def check_upnp():
    try:
        mapper = UpnpPortMapper()
        external_ip = mapper.external_ip()
        return external_ip != ''
    except Exception as ex:
        return False

def check_natpmp():
    try:
        mapper = NatPmpPortMapper()
        external_ip = mapper.external_ip()
        return external_ip != ''
    except Exception as ex:
        return False


def pytest_generate_tests(metafunc):
    if 'mapper' in metafunc.fixturenames:
        ids = []
        mappers = []
        if check_upnp():
            ids.append('UpnpPortMapper')
            mappers.append(UpnpPortMapper())
        if check_natpmp():
            ids.append('NatPmpPortMapper')
            mappers.append(NatPmpPortMapper())
        metafunc.parametrize('mapper', mappers, ids=ids)

def test_external_ip(mapper):
    external_ip = mapper.external_ip()
    assert external_ip is not None

def test_add_mapping_simple(http_server, mapper):
    external_port = mapper.add_mapping(http_server.port)
    assert external_port is not None
    external_ip = mapper.external_ip()
    response = wait_http(external_ip, external_port, 200, timeout=1)
    assert response is not None

def test_add_mapping_twice(http_server, mapper):
    external_port_first = mapper.add_mapping(http_server.port)
    external_port_second = mapper.add_mapping(http_server.port)
    assert external_port_first == external_port_second

def test_remove_mapping(http_server, mapper):
    external_ip = mapper.external_ip()
    local_port = http_server.port
    external_port = mapper.add_mapping(local_port)
    mapper.remove_mapping(local_port, external_port)
    ex = wait_http_cant_connect(external_ip, external_port, timeout=10)
    assert ex is not None
