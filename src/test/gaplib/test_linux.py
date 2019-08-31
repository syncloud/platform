from syncloud_platform.gaplib import linux


def test_is_ip_piblic():
    assert linux.is_ip_public("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
