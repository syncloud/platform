import syncloud_platform.importlib

from syncloud_platform.auth.ldapauth import to_ldap_dc


def test_to_ldap_dc():
    assert to_ldap_dc('user.syncloud.it') == 'dc=user,dc=syncloud,dc=it'