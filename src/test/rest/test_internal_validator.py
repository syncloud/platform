import pytest 

from syncloud_platform.rest.internal_validator import InternalValidator
from syncloudlib.error import PassthroughJsonError


def test_to_json():
    validator = InternalValidator()
    validator.add_parameter_message('login', 'login is empty')
    validator.add_parameter_message('password', 'is too short')
    validator.add_parameter_message('password', 'has no special symbol')
    actual = validator.to_json()
    print(actual)
    expected = '{"parameters_messages": [{"parameter": "login", "messages": ["login is empty"]}, {"parameter": "password", "messages": ["is too short", "has no special symbol"]}]}'
    print(expected)
    assert expected == actual

def test_validate_good_credentials():
    validator = InternalValidator()
    validator.validate('username', 'password123')


def test_validate_short_credentials():
    validator = InternalValidator()
    with pytest.raises(PassthroughJsonError):
        validator.validate('u', 'p')
        