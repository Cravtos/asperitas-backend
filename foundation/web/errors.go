package web

import (
	"github.com/pkg/errors"
)

// Shutdown is a type used to help with the graceful termination of the service.
type Shutdown struct {
	Message string
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful Shutdown.
func NewShutdownError(message string) error {
	return &Shutdown{message}
}

// Error is the implementation of the error interface.
func (s *Shutdown) Error() string {
	return s.Message
}

// IsShutdown checks to see if the Shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*Shutdown); ok {
		return true
	}
	return false
}
