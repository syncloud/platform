import pytest

import responses
from convertible import reformat
from syncloud_platform.board import config
from syncloud_platform.board import footprint
from syncloud_platform.board import id

from syncloud_platform.insider.redirect_service import RedirectService
from syncloud_platform.insider.port_drill import PortDrill
from syncloud_platform.insider.config import Port
from test.insider.helpers import get_port_config, get_user_platform_config

from syncloud_app.main import PassthroughJsonError

def test_version():
    return 'test'

def assertSingleRequest(expected_request):
    expected_request = reformat(expected_request)
    assert 1 == len(responses.calls)
    the_call = responses.calls[0]
    assert expected_request == the_call.request.body


class FakePortDrillProvider:

    def __init__(self, port_drill):
        self.port_drill = port_drill

    def get_drill(self, external_access):
        return self.port_drill


@responses.activate
def test_sync_success():
    port_config = get_port_config([Port(80, 80, 'TCP'), Port(81, 81, 'TCP')])
    port_drill = PortDrill(port_config, MockPortMapper(external_ip='192.167.44.52'), MockPortProber())

    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=200,
                  body="{'message': 'Domain was updated'}",
                  content_type="application/json")

    # platform_config = get_platform_config()

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    user_platform_config.update_domain('boris', 'some_update_token')
    dns = RedirectService(MockNetwork('127.0.0.1'), user_platform_config, test_version)
    dns.sync(port_drill, 'some_update_token', 'http', True, 'TCP')

    expected_request = '''
{
    "web_local_port": 80,
    "web_port": 80,
    "web_protocol": "http",
    "ip": "192.167.44.52",
    "local_ip": "127.0.0.1",
    "map_local_address": false,
    "platform_version": "test",
    "token": "some_update_token"
}
'''
    assertSingleRequest(expected_request)


@responses.activate
def test_sync_server_side_client_ip():
    port_config = get_port_config([Port(80, 80, 'TCP'), Port(81, 81, 'TCP')])
    port_drill = PortDrill(port_config, MockPortMapper(external_ip='10.1.1.1'), MockPortProber())

    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=200,
                  body="{'message': 'Domain was updated'}",
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    user_platform_config.update_domain('boris', 'some_update_token')
    dns = RedirectService(MockNetwork('127.0.0.1'), user_platform_config, test_version)
    dns.sync(port_drill, 'some_update_token', 'http', False, 'TCP')

    expected_request = '''
{
    "web_local_port": 80,
    "web_port": 80,
    "web_protocol": "http",
    "platform_version": "test",
    "token": "some_update_token",
    "map_local_address": true,
    "local_ip": "127.0.0.1"
}
'''
    assertSingleRequest(expected_request)


@responses.activate
def test_sync_server_error():
    port_config = get_port_config([Port(80, 10000, 'TCP')])
    port_drill = PortDrill(port_config, MockPortMapper(external_ip='192.167.44.52'), MockPortProber())

    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=400,
                  body='{"message": "Unknown update token"}',
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    user_platform_config.update_domain('boris', 'some_update_token')
    dns = RedirectService(MockNetwork('127.0.0.1'), user_platform_config, test_version)

    with pytest.raises(PassthroughJsonError) as context:
        dns.sync(port_drill, 'some_update_token', 'http', False, 'TCP')

    assert context.value.message == "Unknown update token"


@responses.activate
def test_link_success():
    config.footprints.append(('my-PC', footprint.footprint()))
    config.titles['my-PC'] = 'My PC'
    device_id = id.id()

    responses.add(responses.POST,
                  "http://api.domain.com/domain/acquire",
                  status=200,
                  body='{"user_domain": "boris", "update_token": "some_update_token"}',
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    dns = RedirectService(MockNetwork('127.0.0.1'), user_platform_config, test_version)
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

    assert result.user_domain == "boris"
    assert result.update_token == "some_update_token"


@responses.activate
def test_link_server_error():
    config.footprints.append(('my-PC', footprint.footprint()))
    config.titles['my-PC'] = 'My PC'

    responses.add(responses.POST,
                  "http://api.domain.com/domain/acquire",
                  status=403,
                  body='{"message": "Authentication failed"}',
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    dns = RedirectService(MockNetwork('127.0.0.1'), user_platform_config, test_version)

    with pytest.raises(PassthroughJsonError) as context:
        result = dns.acquire('boris@mail.com', 'pass1234', 'boris')

    assert context.value.message == "Authentication failed"

    assert user_platform_config.get_user_domain() is None


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


class MockNetwork:

    def __init__(self, ip):
        self.ip = ip

    def local_ip(self):
        return self.ip
