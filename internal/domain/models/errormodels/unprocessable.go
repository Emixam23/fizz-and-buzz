package errormodels

import "encoding/json"

// UnprocessableError is a wrapped error implementation
type UnprocessableError struct {
	error
}

type unprocessableError struct {
	Err string `json:"err"`
}

// NewUnprocessableError creates and initializes an unprocessable error
func NewUnprocessableError(err error) *UnprocessableError {
	return &UnprocessableError{err}
}

// MarshalJSON overrides base implementation
func (e *UnprocessableError) MarshalJSON() ([]byte, error) {
	return json.Marshal(unprocessableError{Err: e.Error()})
}
