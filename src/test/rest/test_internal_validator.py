from syncloud_platform.rest.internal_validator import InternalValidator

def test_errors():
    validator = InternalValidator()
    validator.add_parameter_message('login', 'login is empty')
    validator.add_parameter_message('password', 'is too short')
    validator.add_parameter_message('password', 'has no special symbol')
    print(validator.to_json())
    assert '{"parameters_messages": [{"messages": ["login is empty"], "parameter": "login"}, {"messages": ["is too short", "has no special symbol"], "parameter": "password"}]}' == validator.to_json()