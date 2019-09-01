from syncloud_platform.gaplib import linux


def test_is_ip_piblic():
    assert linux.is_ip_public("8.8.8.8")

def test_is_ip_piblic_false():
    assert not linux.is_ip_public("192.168.0.1")
