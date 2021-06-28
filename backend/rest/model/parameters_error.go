package model

type ParameterError struct {
	ParameterErrors *[]ParameterMessages
}

type ParameterMessages struct {
	Parameter string   `json:"parameter,omitempty"`
	Messages  []string `json:"messages,omitempty"`
}

func (p *ParameterError) Error() string {
	return "There's an error in parameters"
}

func SingleParameterError(parameter string, message string) *ParameterError {
	return &ParameterError{ParameterErrors: &[]ParameterMessages{{
		Parameter: parameter, Messages: []string{message},
	}}}
}
