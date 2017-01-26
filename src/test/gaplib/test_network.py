import socket
from syncloud_platform.gaplib import linux


def test_local_ip():
    assert socket.inet_aton(linux.local_ip())
