package model

type CodedError struct {
	Code  string
	Cause error
}

func (e *CodedError) Error() string {
	return e.Cause.Error()
}

func Coded(code string, err error) *CodedError {
	return &CodedError{Code: code, Cause: err}
}
