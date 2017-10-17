from syncloud_platform.device import create_email

def test_create_email_from_username():
   
    username = 'test@test.com'
    domain = 'example.com'
    
    assert create_email(username, domain) == 'test@test.com'


def test_create_email_from_domain_fallback():
   
    username = 'test'
    domain = 'example.com'
    
    assert create_email(username, domain) == 'test@example.com'
