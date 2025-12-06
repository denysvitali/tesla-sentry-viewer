package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError(t *testing.T) {
	// Test basic error
	err := New("TEST_ERROR", "This is a test error")
	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "This is a test error", err.Message)
	assert.Equal(t, "", err.Details)
	assert.Nil(t, err.Err)
	assert.Equal(t, "TEST_ERROR: This is a test error", err.Error())

	// Test error with details
	err = NewWithDetails("TEST_ERROR", "This is a test error", "Additional details")
	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "This is a test error", err.Message)
	assert.Equal(t, "Additional details", err.Details)
	assert.Nil(t, err.Err)
	assert.Equal(t, "TEST_ERROR: This is a test error (Additional details)", err.Error())

	// Test wrapped error
	originalErr := assert.AnError
	wrappedErr := Wrap(originalErr, "WRAPPED_ERROR", "Wrapped error message")
	assert.Equal(t, "WRAPPED_ERROR", wrappedErr.Code)
	assert.Equal(t, "Wrapped error message", wrappedErr.Message)
	assert.Equal(t, "", wrappedErr.Details)
	assert.Equal(t, assert.AnError, wrappedErr.Err)
	assert.Equal(t, "WRAPPED_ERROR: Wrapped error message", wrappedErr.Error())

	// Test wrapped error with details
	wrappedErr = WrapWithDetails(originalErr, "WRAPPED_ERROR", "Wrapped error message", "More details")
	assert.Equal(t, "WRAPPED_ERROR", wrappedErr.Code)
	assert.Equal(t, "Wrapped error message", wrappedErr.Message)
	assert.Equal(t, "More details", wrappedErr.Details)
	assert.Equal(t, assert.AnError, wrappedErr.Err)
	assert.Equal(t, "WRAPPED_ERROR: Wrapped error message (More details)", wrappedErr.Error())
}

func TestErrorCodes(t *testing.T) {
	// Test that all error codes are defined
	assert.Equal(t, "INVALID_INPUT", ErrInvalidInput)
	assert.Equal(t, "FILE_NOT_FOUND", ErrFileNotFound)
	assert.Equal(t, "DIRECTORY_NOT_FOUND", ErrDirectoryNotFound)
	assert.Equal(t, "PROCESSING_FAILED", ErrProcessingFailed)
	assert.Equal(t, "CONFIG_INVALID", ErrConfigInvalid)
	assert.Equal(t, "PERMISSION_DENIED", ErrPermissionDenied)
	assert.Equal(t, "TIMEOUT", ErrTimeout)
	assert.Equal(t, "NOT_IMPLEMENTED", ErrNotImplemented)
}