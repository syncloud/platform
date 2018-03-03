from syncloud_platform.insider.device_info import construct_url


def test_url_standard_port():
    assert construct_url(443, 'domain.tld') == 'https://domain.tld'


def test_url_non_standard_port():
    assert construct_url(444, 'domain.tld') == 'https://domain.tld:444'
