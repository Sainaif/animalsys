package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements error interface
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

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with message and status code
func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common errors
var (
	ErrBadRequest          = New(http.StatusBadRequest, "Bad request")
	ErrUnauthorized        = New(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "Forbidden")
	ErrNotFound            = New(http.StatusNotFound, "Resource not found")
	ErrConflict            = New(http.StatusConflict, "Resource already exists")
	ErrInternalServer      = New(http.StatusInternalServerError, "Internal server error")
	ErrInvalidCredentials  = New(http.StatusUnauthorized, "Invalid email or password")
	ErrInvalidToken        = New(http.StatusUnauthorized, "Invalid or expired token")
	ErrEmailAlreadyExists  = New(http.StatusConflict, "Email already exists")
	ErrValidationFailed    = New(http.StatusBadRequest, "Validation failed")
)

// NewBadRequest creates a bad request error with custom message
func NewBadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

// NewNotFound creates a not found error with custom message
func NewNotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

// NewUnauthorized creates an unauthorized error with custom message
func NewUnauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

// NewForbidden creates a forbidden error with custom message
func NewForbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

// NewConflict creates a conflict error with custom message
func NewConflict(message string) *AppError {
	return New(http.StatusConflict, message)
}

// NewInternalServer creates an internal server error with custom message
func NewInternalServer(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}
