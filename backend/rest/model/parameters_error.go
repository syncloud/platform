package model

import (
	"fmt"
	"strings"
)

type ParameterError struct {
	ParameterErrors *[]ParameterMessages
}

type ParameterMessages struct {
	Parameter string   `json:"parameter,omitempty"`
	Messages  []string `json:"messages,omitempty"`
}

func (pm *ParameterMessages) Error() string {
	return fmt.Sprintf("%s: %s", pm.Parameter, strings.Join(pm.Messages, ", "))
}

func (p *ParameterError) Error() string {
	var errors []string
	for _, pm := range *p.ParameterErrors {
		errors = append(errors, pm.Error())
	}
	return fmt.Sprintf("There's an error in parameters: %s", strings.Join(errors, "; "))
}

func SingleParameterError(parameter string, message string) *ParameterError {
	return &ParameterError{ParameterErrors: &[]ParameterMessages{{
		Parameter: parameter, Messages: []string{message},
	}}}
}
