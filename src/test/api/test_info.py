from syncloud_platform.api.info import __url


def test_url_no_prototol():
    assert __url(None, 80, 'domain.tld') == 'http://domain.tld'


def test_url_standard_port():
    assert __url('https', 443, 'domain.tld') == 'https://domain.tld'


def test_url_non_standard_port():
    assert __url('https', 444, 'domain.tld') == 'https://domain.tld:444'


def test_url_subdomain():
    assert __url('https', 444, 'domain.tld', 'app') == 'https://app.domain.tld:444'
