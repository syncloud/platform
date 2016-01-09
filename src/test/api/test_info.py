from syncloud_platform.insider.device_info import construct_url


def test_url_no_protocol():
    assert construct_url(None, 80, 'domain.tld') == 'http://domain.tld'


def test_url_standard_port():
    assert construct_url('https', 443, 'domain.tld') == 'https://domain.tld'


def test_url_non_standard_port():
    assert construct_url('https', 444, 'domain.tld') == 'https://domain.tld:444'


def test_url_sub_domain():
    assert construct_url('https', 444, 'domain.tld', 'app') == 'https://app.domain.tld:444'
