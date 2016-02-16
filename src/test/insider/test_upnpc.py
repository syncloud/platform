from syncloud_platform.insider.upnpc import get_our_their_ports, Mapping


def test_get_our_their_ports():
    our_expected = Mapping(1, 'TCP', '127.0.0.1', 80, 'desc1', True, '127.1.1.1', 0)
    their_expected = Mapping(2, 'TCP', '127.0.0.2', 80, 'desc1', True, '127.1.1.1', 0)

    our_ports, their_ports = get_our_their_ports([our_expected, their_expected], '127.0.0.1', 'TCP', 80)

    assert our_ports == [1]
    assert their_ports == [2]
