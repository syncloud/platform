from syncloud_platform.control.tls import certificate_is_valid


def test_certificate_is_valid_60days_no_new_domains():

    assert certificate_is_valid(60, []) is True


def test_certificate_is_valid_60days_new_domains():

    assert certificate_is_valid(60, ['a']) is False


def test_certificate_is_valid_10days_no_new_domains():

    assert certificate_is_valid(10, []) is False
