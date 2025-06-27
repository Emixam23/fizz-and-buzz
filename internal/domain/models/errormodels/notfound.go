package errormodels

import "encoding/json"

// NotFoundError is a wrapped error implementation
type NotFoundError struct {
	error
}

type notFoundError struct {
	Err string `json:"err"`
}

// NewNotFoundError creates and initializes a not found error
func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}

// MarshalJSON overrides base implementation
func (e *NotFoundError) MarshalJSON() ([]byte, error) {
	return json.Marshal(notFoundError{Err: e.Error()})
}
