from syncloud_platform.auth.ldapauth import generate_change_password_cmd


def test_generate_change_password_cmd():
    assert generate_change_password_cmd('123123') == 'echo "root:123123" | chpasswd'
    assert generate_change_password_cmd('123"123') == 'echo "root:123\\"123" | chpasswd'
