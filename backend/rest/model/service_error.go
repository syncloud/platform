package model

import "fmt"

type ServiceError struct {
	InternalError error
	StatusCode    int
}

func (e *ServiceError) Error() string {
	return e.InternalError.Error()
}

func NewServiceErrorWithCode(message string, code int) *ServiceError {
	return &ServiceError{InternalError: fmt.Errorf(message), StatusCode: code}
}
