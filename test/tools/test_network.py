import socket
from syncloud.tools import network


def test_local_ip():
    assert socket.inet_aton(network.local_ip())