import responses

from syncloud_platform.insider.port_prober import PortProber


@responses.activate
def test_private_ip():

    responses.add(responses.GET,
                  "http://api.domain.com/probe/port_v2",
                  status=200,
                  body='{"message": "OK", "device_ip": "0"}',
                  content_type="application/json")

    prober = PortProber('http://api.domain.com', '')
    result, message = prober.probe_port(80, 'http', '0')

    assert len(responses.calls) == 0
    assert result is False


@responses.activate
def test_no_ip():

    responses.add(responses.GET,
                  "http://api.domain.com/probe/port_v2",
                  status=200,
                  body='{"message": "OK", "device_ip": "0"}',
                  content_type="application/json")

    prober = PortProber('http://api.domain.com', '')
    result, message = prober.probe_port(80, 'http', None)

    assert result is True
    assert 'ip=' not in responses.calls[0].request.url


@responses.activate
def test_public_ip():

    responses.add(responses.GET,
                  "http://api.domain.com/probe/port_v2",
                  status=200,
                  body='{"message": "OK", "device_ip": "0"}',
                  content_type="application/json")

    prober = PortProber('http://api.domain.com', '')
    result, message = prober.probe_port(80, 'http', '8.8.8.8')

    assert result is True
    assert 'ip=' in responses.calls[0].request.url
