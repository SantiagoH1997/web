package web

import "github.com/pkg/errors"

// Error represents a generic error.
type Error struct {
	Err        error
	StatusCode int
	Fields     []ErrorField
}

// ErrorResponse represents the error
// to be sent to the client.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []ErrorField `json:"fields,omitempty"`
}

// ErrorField represents a field that failed validation.
type ErrorField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func (re *Error) Error() string {
	return re.Err.Error()
}

// NewError returns an error with an HTTP status code.
func NewError(err error, statusCode int) error {
	return &Error{err, statusCode, nil}
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &Error{err, status, nil}
}

// Error is the implementation of the error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// IsShutdown checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
