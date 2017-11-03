from syncloud_platform.device import parse_username

def test_parse_username_from_username():
   
    username = 'test@test.com'
    domain = 'example.com'
    
    name, email = parse_username(username, domain)
    assert name == 'test'
    assert email == 'test@test.com'


def test_parse_username_from_domain_fallback():
   
    username = 'test'
    domain = 'example.com'
    
    name, email = parse_username(username, domain)
    assert name == 'test'
    assert email == 'test@example.com'
