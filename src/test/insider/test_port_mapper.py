from miniupnpc import UPnP

import pytest
import time

from syncloud_platform.insider.natpmpc import NatPmpPortMapper
from syncloud_platform.insider.upnpc import UpnpPortMapper
from syncloud_platform.insider.port_drill import provide_mapper
from test.insider.http import SomeHttpServer, wait_http, wait_http_cant_connect


@pytest.fixture(scope="module")
def http_server(request):
    server = SomeHttpServer(18088)
    server.start()

    def fin():
        server.stop()
    request.addfinalizer(fin)
    return server


ids = []
mappers = []

mapper = provide_mapper(NatPmpPortMapper(), UpnpPortMapper(UPnP()))
if mapper is not None:
    ids.append(mapper.name())
    mappers.append(mapper)


@pytest.mark.skip(reason="Port mapping is very unstable on build server")
@pytest.mark.parametrize("mapper", mappers, ids=ids)
def test_external_ip(mapper):
    external_ip = mapper.external_ip()
    assert external_ip is not None


@pytest.mark.skip(reason="Port mapping is very unstable on build server")
@pytest.mark.parametrize("mapper", mappers, ids=ids)
def test_add_mapping_simple(http_server, mapper):
    external_port = mapper.add_mapping(http_server.port, http_server.port, 'TCP')
    assert external_port is not None
    external_ip = mapper.external_ip()
    response = wait_http(external_ip, external_port, 200, timeout=2)
    print(response)
    assert response is not None


@pytest.mark.skip(reason="Port mapping is very unstable on build server")
@pytest.mark.parametrize("mapper", mappers, ids=ids)
def test_add_mapping_twice(http_server, mapper):
    external_port_first = mapper.add_mapping(http_server.port, http_server.port, 'TCP')
    external_port_second = mapper.add_mapping(http_server.port, http_server.port, 'TCP')
    assert external_port_first == external_port_second


@pytest.mark.skip(reason="Port mapping is very unstable on build server")
@pytest.mark.parametrize("mapper", mappers, ids=ids)
def test_remove_mapping(http_server, mapper):
    external_ip = mapper.external_ip()
    local_port = http_server.port
    external_port = mapper.add_mapping(local_port, http_server.port, 'TCP')
    mapper.remove_mapping(local_port, external_port, 'TCP')
    ex = wait_http_cant_connect(external_ip, external_port, timeout=5)
    assert ex is not None
