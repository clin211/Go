package errors

import (
	"fmt"
)

// AppError is a custom error type for the application
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// ValidationError creates a new validation error
func ValidationError(message string, err error) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

// NotFoundError creates a new not found error
func NotFoundError(message string, err error) *AppError {
	return &AppError{
		Code:    404,
		Message: message,
		Err:     err,
	}
}

// InternalError creates a new internal server error
func InternalError(message string, err error) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Code == 400
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Code == 404
}

// IsInternalError checks if an error is an internal server error
func IsInternalError(err error) bool {
	appErr, ok := err.(*AppError)
	return ok && appErr.Code == 500
}
