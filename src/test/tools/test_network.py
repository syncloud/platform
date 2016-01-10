
import socket
from syncloud_platform.tools.network import Network


def test_local_ip():
    assert socket.inet_aton(Network().local_ip())
