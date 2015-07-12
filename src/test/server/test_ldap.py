from syncloud_platform.server.auth import to_ldap_dc


def test_to_ldap_dc():
    assert to_ldap_dc('user.syncloud.it') == 'dc=user,dc=syncloud,dc=it'