from syncloud_platform.device import create_email

def test_create_email_from_username():
   
    username = 'test@test.com'
    domain = 'example.com'
    
    name, email = create_email(username, domain)
    assert name == 'test'
    assert name == 'test@test.com'


def test_create_email_from_domain_fallback():
   
    username = 'test'
    domain = 'example.com'
    
    name, email = create_email(username, domain)
    assert name == 'test'
    assert email == 'test@example.com'
