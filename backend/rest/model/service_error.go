package model

import (
	"net/http"
)

type ServiceError struct {
	InternalError error
	StatusCode    int
}

func (e *ServiceError) Error() string {
	return e.InternalError.Error()
}

func BadRequest(err error) *ServiceError {
	return &ServiceError{InternalError: err, StatusCode: http.StatusBadRequest}
}
