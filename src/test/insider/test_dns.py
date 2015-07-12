import logging
import pytest

import responses
from syncloud_app import logger
from convertible import reformat
from syncloud_platform.tools import config
from syncloud_platform.tools import footprint
from syncloud_platform.tools import id

from syncloud_platform.insider.dns import Dns
from syncloud_platform.insider.port_mapper import PortMapper, MockPortMapper
from syncloud_platform.insider.config import Port, Domain, Service
from test.insider.helpers import get_port_config, get_domain_config, get_service_config, get_insider_config

from syncloud_app.main import PassthroughJsonError

logger.init(level=logging.DEBUG, console=True)


def assertSingleRequest(expected_request):
    expected_request = reformat(expected_request)
    assert 1 == len(responses.calls)
    the_call = responses.calls[0]
    assert expected_request == the_call.request.body

@responses.activate
def test_sync_success():
    service_config = get_service_config([
        Service("ownCloud", "http", "_http._tcp", 80, url="owncloud"),
        Service("SSH", "https", "_http._tcp", 81, url=None)
    ])
    port_config = get_port_config([Port(80, 80), Port(81, 81)])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='192.167.44.52'))

    domain_config = get_domain_config(Domain('boris', 'some_update_token'))

    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=200,
                  body="{'message': 'Domain was updated'}",
                  content_type="application/json")

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')
    dns.sync()

    expected_request = '''
{
    "services": [
        {"name": "ownCloud", "protocol": "http", "type": "_http._tcp", "port": 80, "local_port": 80, "url": "owncloud"},
        {"name": "SSH", "protocol": "https", "type": "_http._tcp", "port": 81, "local_port": 81, "url": null}
    ],
    "ip": "192.167.44.52",
    "local_ip": "127.0.0.1",
    "token": "some_update_token"
}
'''
    assertSingleRequest(expected_request)

@responses.activate
def test_sync_server_side_client_ip():
    service_config = get_service_config([
        Service("ownCloud", "http", "_http._tcp", 80, url="owncloud"),
        Service("SSH", "https", "_http._tcp", 81, url=None)
    ])
    port_config = get_port_config([Port(80, 80), Port(81, 81)])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='10.1.1.1'))

    domain_config = get_domain_config(Domain('boris', 'some_update_token'))

    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=200,
                  body="{'message': 'Domain was updated'}",
                  content_type="application/json")

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')
    dns.sync()

    expected_request = '''
{
    "services": [
        {"name": "ownCloud", "protocol": "http", "type": "_http._tcp", "port": 80, "local_port": 80, "url": "owncloud"},
        {"name": "SSH", "protocol": "https", "type": "_http._tcp", "port": 81, "local_port": 81, "url": null}
    ],
    "token": "some_update_token",
    "local_ip": "127.0.0.1"
}
'''
    assertSingleRequest(expected_request)


@responses.activate
def test_sync_server_error():
    service_config = get_service_config([Service("ownCloud", "http", "_http._tcp", 80, url="owncloud")])
    port_config = get_port_config([Port(80, 10000)])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='192.167.44.52'))

    domain_config = get_domain_config(Domain('boris', 'some_update_token'))

    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=400,
                  body='{"message": "Unknown update token"}',
                  content_type="application/json")

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')

    with pytest.raises(PassthroughJsonError) as context:
        dns.sync()

    assert context.value.message == "Unknown update token"

@responses.activate
def test_link_success():
    config.footprints.append(('my-PC', footprint.footprint()))
    config.titles['my-PC'] = 'My PC'
    device_id = id.id()

    domain_config = get_domain_config(None)

    responses.add(responses.POST,
                  "http://api.domain.com/domain/acquire",
                  status=200,
                  body='{"user_domain": "boris", "update_token": "some_update_token"}',
                  content_type="application/json")

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config=None, port_mapper=None, local_ip='127.0.0.1')
    result = dns.acquire('boris@mail.com', 'pass1234', 'boris')

    assert result is not None
    assert result.user_domain == "boris"
    assert result.update_token == "some_update_token"

    expected_request_data = {
        "password": "pass1234",
        "email": "boris@mail.com",
        "user_domain": "boris",
        'device_mac_address': device_id.mac_address,
        'device_name': device_id.name,
        'device_title': device_id.title,
    }
    # Need to assert all passed POST parameters
    # self.assertSingleRequest(convertible.to_json(expected_request_data))

    domain = domain_config.load()
    assert domain is not None
    assert domain.user_domain == "boris"
    assert domain.update_token == "some_update_token"

@responses.activate
def test_link_server_error():
    config.footprints.append(('my-PC', footprint.footprint()))
    config.titles['my-PC'] = 'My PC'

    domain_config = get_domain_config(None)

    responses.add(responses.POST,
                  "http://api.domain.com/domain/acquire",
                  status=403,
                  body='{"message": "Authentication failed"}',
                  content_type="application/json")

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config=None, port_mapper=None, local_ip='127.0.0.1')
    with pytest.raises(PassthroughJsonError) as context:
        result = dns.acquire('boris@mail.com', 'pass1234', 'boris')

    assert context.value.message == "Authentication failed"

    domain = domain_config.load()
    assert domain is None

def test_add_service():
    service_config = get_service_config([])
    port_config = get_port_config([])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='192.167.44.52'))

    domain_config = get_domain_config(None)

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')
    dns.add_service("ownCloud", "http", "_http._tcp", 80, url="owncloud")

    services = service_config.load()
    assert 1 == len(services)
    service = services[0]
    assert "ownCloud" == service.name
    assert "_http._tcp" == service.type
    assert 80 == service.port
    assert "owncloud" == service.url

    mappings = port_config.load()
    assert 1 == len(mappings)
    mapping = mappings[0]
    assert 80 == mapping.local_port

def test_get_service():
    service_config = get_service_config([])
    port_config = get_port_config([])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='192.167.44.52'))

    domain_config = get_domain_config(None)

    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')
    dns.add_service("ownCloud", "http", "_http._tcp", 80, url="owncloud")

    service = dns.get_service("ownCloud")

    assert service is not None
    assert "ownCloud" == service.name
    assert "_http._tcp" == service.type
    assert 80 == service.port
    assert "owncloud" == service.url

def test_get_not_existing_service():
    service_config = get_service_config([])
    port_config = get_port_config([])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='192.167.44.52'))

    domain_config = get_domain_config(None)
    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')

    service = dns.get_service("ownCloud")

    assert service is None

def test_endpoints():
    service_config = get_service_config([
        Service("ownCloud", "http", "_http._tcp", 80, url="owncloud"),
        Service("SSH", "https", "_http._tcp", 81, url=None)
    ])
    port_config = get_port_config([Port(80, 8080), Port(81, 8181)])
    port_mapper = PortMapper(port_config, lambda: MockPortMapper(external_ip='10.1.1.1'))

    domain_config = get_domain_config(Domain('boris', 'some_update_token'))
    insider_config = get_insider_config('domain.com', 'http://api.domain.com')
    dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')

    endpoints = dns.endpoints()
    assert len(endpoints) == 2
    assert endpoints[0].service.name == 'ownCloud'
    assert endpoints[0].external_port == 8080
    assert endpoints[1].service.name == 'SSH'
    assert endpoints[1].external_port == 8181