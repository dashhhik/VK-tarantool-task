package core

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound    = errors.New("resource not found")
	ErrForbidden   = errors.New("access forbidden")
	ErrKeyNotFound = errors.New("key not found")
)

// CustomError allows you to create errors with a code and message.
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// NewCustomError creates a new CustomError with a code and message.
func NewCustomError(code int, message string) error {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// WrapError allows wrapping another error with additional context.
func WrapError(context string, err error) error {
	return fmt.Errorf("%s: %w", context, err)
}
