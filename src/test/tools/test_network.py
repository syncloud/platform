import syncloud_platform.importlib

import socket
from syncloud_platform.tools import network


def test_local_ip():
    assert socket.inet_aton(network.local_ip())