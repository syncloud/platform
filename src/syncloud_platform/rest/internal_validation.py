import json

class ParameterMessages:

    def __init__(self, parameter):
        self.parameter = parameter
        self.messages = []
        
    def add_message(self, message):
        self.messages.append(message)
  
class InternalValidator:
    def __init__(self):
        self.parameter_messages = { }
        
    def add_parameter_message(self, parameter, message):

        if parameter not in self.parameter_messages:
            parameter_messages = ParameterMessages(parameter)
            self.parameter_messages[parameter] = parameter_messages
            
        self.parameter_messages[parameter].add_message(message)
    
    def to_json(self):
        parameters_messages = [ v.__dict__ for k,v in self.parameter_messages.iteritems() ]
        return json.dumps({ 'parameters_messages': parameters_messages })
            
