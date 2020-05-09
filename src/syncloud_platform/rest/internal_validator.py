import json

from syncloud_platform.rest.model.parameter_messages import ParameterMessages
from syncloudlib.error import PassthroughJsonError
 
 
class InternalValidator:
    def __init__(self):
        self.parameter_messages = { }
    
    def validate(self, device_username, device_password):
        if len(device_username) < 3:
            self.add_parameter_message('device_username', 'less than 3 characters')
        
        if '!' in device_username:
            self.add_parameter_message('device_username', 'contains ! symbol')
   
        if len(device_password) < 7:
            self.add_parameter_message('device_password', 'less than 7 characters')
            
        if len(self.parameter_messages) != 0:
            raise PassthroughJsonError('validation errors', self.to_json())
        
    def add_parameter_message(self, parameter, message):

        if parameter not in self.parameter_messages:
            parameter_messages = ParameterMessages(parameter)
            self.parameter_messages[parameter] = parameter_messages
            
        self.parameter_messages[parameter].add_message(message)
    
    def to_json(self):
        parameters_messages = [ v.__dict__ for k,v in self.parameter_messages.items() ]
        return json.dumps({ 'parameters_messages': parameters_messages })
            
   
