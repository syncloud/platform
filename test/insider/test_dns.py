import logging
import unittest

import responses
from mock import MagicMock
from syncloud.app import logger
from convertible import reformat
from syncloud.tools import config
from syncloud.tools import footprint
from syncloud.tools import id

from syncloud.insider.dns import Dns
from syncloud.insider.port_mapper import PortMapper
from syncloud.insider.config import Port, Domain, Service
from test.insider.helpers import get_port_config, get_domain_config, get_service_config, get_insider_config

from syncloud.app.main import PassthroughJsonError

logger.init(level=logging.DEBUG, console=True)


def get_upnpc(external_ip):
    upnpc = MagicMock()
    upnpc.external_ip = MagicMock(return_value=external_ip)
    upnpc.port_open_on_router = MagicMock(return_value=False)
    return upnpc


def get_port_mapper(port_config, upnpc):
    return PortMapper(port_config, upnpc)


class TestDns(unittest.TestCase):

    def assertSingleRequest(self, expected_request):
        expected_request = reformat(expected_request)
        self.assertEquals(1, len(responses.calls))
        the_call = responses.calls[0]
        self.assertEquals(expected_request, the_call.request.body)

    @responses.activate
    def test_sync_success(self):
        service_config = get_service_config([
            Service("ownCloud", "http", "_http._tcp", 80, url="owncloud"),
            Service("SSH", "https", "_http._tcp", 81, url=None)
        ])
        port_config = get_port_config([Port(80, 80), Port(81, 81)])
        upnpc = get_upnpc(external_ip='192.167.44.52')
        upnpc.mapped_external_ports = MagicMock(side_effect=[[], [80]])

        port_mapper = get_port_mapper(port_config, upnpc)
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
        self.assertSingleRequest(expected_request)

    @responses.activate
    def test_sync_server_side_client_ip(self):
        service_config = get_service_config([
            Service("ownCloud", "http", "_http._tcp", 80, url="owncloud"),
            Service("SSH", "https", "_http._tcp", 81, url=None)
        ])
        port_config = get_port_config([Port(80, 80), Port(81, 81)])
        upnpc = get_upnpc(external_ip='10.1.1.1')
        upnpc.mapped_external_ports = MagicMock(side_effect=[[], [80]])

        port_mapper = get_port_mapper(port_config, upnpc)
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
        self.assertSingleRequest(expected_request)


    @responses.activate
    def test_sync_server_error(self):
        service_config = get_service_config([Service("ownCloud", "http", "_http._tcp", 80, url="owncloud")])
        port_config = get_port_config([Port(80, 10000)])
        upnpc = get_upnpc(external_ip='192.167.44.52')
        port_mapper = get_port_mapper(port_config, upnpc)
        domain_config = get_domain_config(Domain('boris', 'some_update_token'))

        responses.add(responses.POST,
                      "http://api.domain.com/domain/update",
                      status=400,
                      body='{"message": "Unknown update token"}',
                      content_type="application/json")

        insider_config = get_insider_config('domain.com', 'http://api.domain.com')
        dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')

        with self.assertRaises(PassthroughJsonError) as context:
            dns.sync()

        self.assertEquals(context.exception.message, "Unknown update token")

    @responses.activate
    def test_link_success(self):
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

        self.assertIsNotNone(result)
        self.assertEquals(result.user_domain, "boris")
        self.assertEquals(result.update_token, "some_update_token")

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
        self.assertIsNotNone(domain)
        self.assertEquals(domain.user_domain, "boris")
        self.assertEquals(domain.update_token, "some_update_token")

    @responses.activate
    def test_link_server_error(self):
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
        with self.assertRaises(PassthroughJsonError) as context:
            result = dns.acquire('boris@mail.com', 'pass1234', 'boris')

        self.assertEquals(context.exception.message, "Authentication failed")

        domain = domain_config.load()
        self.assertIsNone(domain)

    def test_add_service(self):
        service_config = get_service_config([])
        port_config = get_port_config([])
        upnpc = get_upnpc(external_ip='192.167.44.52')
        port_mapper = get_port_mapper(port_config, upnpc)

        domain_config = get_domain_config(None)

        insider_config = get_insider_config('domain.com', 'http://api.domain.com')
        dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')
        dns.add_service("ownCloud", "http", "_http._tcp", 80, url="owncloud")

        services = service_config.load()
        self.assertEquals(1, len(services))
        service = services[0]
        self.assertEquals("ownCloud", service.name)
        self.assertEquals("_http._tcp", service.type)
        self.assertEquals(80, service.port)
        self.assertEquals("owncloud", service.url)

        mappings = port_config.load()
        self.assertEquals(1, len(mappings))
        mapping = mappings[0]
        self.assertEquals(80, mapping.local_port)

    def test_get_service(self):
        service_config = get_service_config([])
        port_config = get_port_config([])
        upnpc = get_upnpc(external_ip='192.167.44.52')
        port_mapper = get_port_mapper(port_config, upnpc)

        domain_config = get_domain_config(None)

        insider_config = get_insider_config('domain.com', 'http://api.domain.com')
        dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')
        dns.add_service("ownCloud", "http", "_http._tcp", 80, url="owncloud")

        service = dns.get_service("ownCloud")

        self.assertIsNotNone(service)
        self.assertEquals("ownCloud", service.name)
        self.assertEquals("_http._tcp", service.type)
        self.assertEquals(80, service.port)
        self.assertEquals("owncloud", service.url)

    def test_get_not_existing_service(self):
        service_config = get_service_config([])
        port_config = get_port_config([])
        upnpc = get_upnpc(external_ip='192.167.44.52')
        port_mapper = get_port_mapper(port_config, upnpc)

        domain_config = get_domain_config(None)
        insider_config = get_insider_config('domain.com', 'http://api.domain.com')
        dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')

        service = dns.get_service("ownCloud")

        self.assertIsNone(service)

    def test_endpoints(self):
        service_config = get_service_config([
            Service("ownCloud", "http", "_http._tcp", 80, url="owncloud"),
            Service("SSH", "https", "_http._tcp", 81, url=None)
        ])
        port_config = get_port_config([Port(80, 8080), Port(81, 8181)])
        upnpc = get_upnpc(external_ip='10.1.1.1')
        # upnpc.mapped_external_ports = MagicMock(side_effect=[[], [80]])

        port_mapper = get_port_mapper(port_config, upnpc)
        domain_config = get_domain_config(Domain('boris', 'some_update_token'))
        insider_config = get_insider_config('domain.com', 'http://api.domain.com')
        dns = Dns(insider_config, domain_config, service_config, port_mapper, '127.0.0.1')

        endpoints = dns.endpoints()
        assert len(endpoints) == 2
        assert endpoints[0].service.name == 'ownCloud'
        assert endpoints[0].external_port == 8080
        assert endpoints[1].service.name == 'SSH'
        assert endpoints[1].external_port == 8181