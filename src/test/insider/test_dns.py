import pytest

import responses
from convertible import reformat
from syncloud_platform.board import id

from syncloud_platform.gaplib import linux

def mock_local_ip(ip):
    def fake_local_ip():
        return ip
    return fake_local_ip

linux.local_ip = mock_local_ip('127.0.0.1')

from syncloud_platform.insider.redirect_service import RedirectService
from test.insider.helpers import get_user_platform_config

from syncloud_app.main import PassthroughJsonError

class TestVersions:
    def __init__(self, version):
        self.version = version
    def platform_version(self):
        return self.version


def assertSingleRequest(expected_request):
    expected_request = reformat(expected_request)
    assert 1 == len(responses.calls)
    the_call = responses.calls[0]
    assert expected_request == the_call.request.body


@responses.activate
def test_sync_success():
    
    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=200,
                  body="{'message': 'Domain was updated'}",
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    user_platform_config.update_domain('boris', 'some_update_token')
    dns = RedirectService(user_platform_config, TestVersions('test'))
    dns.sync('192.167.44.52', 80, 80, 'http', 'some_update_token', True)

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
    
    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=200,
                  body="{'message': 'Domain was updated'}",
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    user_platform_config.update_domain('boris', 'some_update_token')
    dns = RedirectService(user_platform_config, TestVersions('test'))
    dns.sync('192.167.44.52', 80, 80, 'http', 'some_update_token', False)

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
    
    responses.add(responses.POST,
                  "http://api.domain.com/domain/update",
                  status=400,
                  body='{"message": "Unknown update token"}',
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    user_platform_config.update_domain('boris', 'some_update_token')
    dns = RedirectService(user_platform_config, TestVersions('test'))

    with pytest.raises(PassthroughJsonError) as context:
        dns.sync('192.167.44.52', 80, 80, 'http', 'some_update_token', False) 

    assert context.value.message == "Unknown update token"


@responses.activate
def test_link_success():
    device_id = id.id()

    responses.add(responses.POST,
                  "http://api.domain.com/domain/acquire",
                  status=200,
                  body='{"user_domain": "boris", "update_token": "some_update_token"}',
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    dns = RedirectService(user_platform_config, TestVersions('test'))
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
    responses.add(responses.POST,
                  "http://api.domain.com/domain/acquire",
                  status=403,
                  body='{"message": "Authentication failed"}',
                  content_type="application/json")

    user_platform_config = get_user_platform_config()
    user_platform_config.update_redirect('domain.com', 'http://api.domain.com')
    dns = RedirectService(user_platform_config, TestVersions('test'))

    with pytest.raises(PassthroughJsonError) as context:
        result = dns.acquire('boris@mail.com', 'pass1234', 'boris')

    assert context.value.message == "Authentication failed"

    assert user_platform_config.get_user_domain() is None
