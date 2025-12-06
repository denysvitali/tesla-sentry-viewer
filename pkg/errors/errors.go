package errors

import (
	"fmt"
)

type AppError struct {
	Code    string
	Message string
	Details string
	Err     error
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewWithDetails(code, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func WrapWithDetails(err error, code, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// Common error codes
const (
	ErrInvalidInput      = "INVALID_INPUT"
	ErrFileNotFound      = "FILE_NOT_FOUND"
	ErrDirectoryNotFound = "DIRECTORY_NOT_FOUND"
	ErrProcessingFailed  = "PROCESSING_FAILED"
	ErrConfigInvalid     = "CONFIG_INVALID"
	ErrPermissionDenied  = "PERMISSION_DENIED"
	ErrTimeout           = "TIMEOUT"
	ErrNotImplemented     = "NOT_IMPLEMENTED"
)