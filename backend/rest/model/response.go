package model

type Response struct {
	Success            bool                 `json:"success"`
	Message            string               `json:"message,omitempty"`
	Code               string               `json:"code,omitempty"`
	Data               *interface{}         `json:"data,omitempty"`
	ParametersMessages *[]ParameterMessages `json:"parameters_messages,omitempty"`
}
