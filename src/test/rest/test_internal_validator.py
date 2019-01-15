import pytest 

from syncloud_platform.rest.internal_validator import InternalValidator
from syncloudlib.error import PassthroughJsonError


def test_to_json():
    validator = InternalValidator()
    validator.add_parameter_message('login', 'login is empty')
    validator.add_parameter_message('password', 'is too short')
    validator.add_parameter_message('password', 'has no special symbol')
    print(validator.to_json())
    assert '{"parameters_messages": [{"messages": ["login is empty"], "parameter": "login"}, {"messages": ["is too short", "has no special symbol"], "parameter": "password"}]}' == validator.to_json()

def test_validate_good_credentials():
    validator = InternalValidator()
    validator.validate('username', 'password123')


def test_validate_short_credentials():
    validator = InternalValidator()
    with pytest.raises(PassthroughJsonError):
        validator.validate('u', 'p')
        