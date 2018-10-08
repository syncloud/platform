from syncloud_platform.auth.ldapauth import generate_change_password_cmd
from syncloud_platform.auth.ldapauth import to_ldap_dc


def test_generate_change_password_cmd():
    assert generate_change_password_cmd('123123') == 'echo "root:123123" | chpasswd'
    assert generate_change_password_cmd('123"123') == 'echo "root:123\\"123" | chpasswd'
    assert generate_change_password_cmd('123$123') == 'echo "root:123\$123" | chpasswd'


def test_to_ldap_dc():
    assert to_ldap_dc('user.syncloud.it') == 'dc=user,dc=syncloud,dc=it'
