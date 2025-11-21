// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error implements the error interface.
type Error struct {
	Id     string `json:"id"`
	Code   int32  `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(id, detail string, code int32) error {
	return &Error{
		Id:     id,
		Code:   code,
		Detail: detail,
		Status: http.StatusText(int(code)),
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

// BadRequest generates a 400 error.
func BadRequest(a ...interface{}) error {
	return &Error{
		Id:     "400",
		Code:   400,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(400),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(a ...interface{}) error {
	return &Error{
		Id:     "401",
		Code:   401,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(401),
	}
}

// Forbidden generates a 403 error.
func Forbidden(a ...interface{}) error {
	return &Error{
		Id:     "403",
		Code:   403,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(403),
	}
}

// NotFound generates a 404 error.
func NotFound(a ...interface{}) error {
	return &Error{
		Id:     "404",
		Code:   404,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(404),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(a ...interface{}) error {
	return &Error{
		Id:     "405",
		Code:   405,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(405),
	}
}

// Timeout generates a 408 error.
func Timeout(a ...interface{}) error {
	return &Error{
		Id:     "408",
		Code:   408,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(408),
	}
}

// Conflict generates a 409 error.
func Conflict(a ...interface{}) error {
	return &Error{
		Id:     "409",
		Code:   409,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(409),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(a ...interface{}) error {
	return &Error{
		Id:     "500",
		Code:   500,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(500),
	}
}

// ServiceUnavailable generates a 503 error.
func ServiceUnavailable(a ...interface{}) error {
	return &Error{
		Id:     "503",
		Code:   503,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(503),
	}
}

// GatewayTimeout generates a 504 error.
func GatewayTimeout(a ...interface{}) error {
	return &Error{
		Id:     "504",
		Code:   504,
		Detail: fmt.Sprintf("%s", a...),
		Status: http.StatusText(504),
	}
}

// ============================================================================
// Domain Specific Errors - Collateral Service
// ============================================================================

// Collateral domain errors
var (
	// ErrCollateralNotFound represents when collateral is not found
	ErrCollateralNotFound = NotFound("collateral not found")

	// ErrCollateralExists represents when collateral already exists
	ErrCollateralExists = Conflict("collateral already exists")

	// ErrInvalidCollateralInput represents invalid collateral input
	ErrInvalidCollateralInput = BadRequest("invalid collateral input")

	// ErrCollateralSyncFailed represents collateral synchronization failure
	ErrCollateralSyncFailed = ServiceUnavailable("collateral synchronization failed")

	// ErrDatabaseOperation represents database operation failures
	ErrDatabaseOperation = InternalServerError("database operation failed")

	// ErrExternalService represents external service failures
	ErrExternalService = ServiceUnavailable("external service unavailable")
)

// ============================================================================
// Error Helper Functions
// ============================================================================

// Is checks if the error is of a specific type
func Is(err error, target error) bool {
	if err == nil || target == nil {
		return false
	}

	// For our Error type, compare the error ID
	if e, ok := err.(*Error); ok {
		if t, ok := target.(*Error); ok {
			return e.Id == t.Id
		}
	}

	// Fallback to simple string comparison
	return err.Error() == target.Error()
}

// Wrap creates a new error with additional context while preserving the original error
func Wrap(originalError error, message string) error {
	if originalError == nil {
		return New("WRAPPED_ERROR", message, 500)
	}

	// If it's already our Error type, enhance it
	if e, ok := originalError.(*Error); ok {
		return &Error{
			Id:     e.Id,
			Code:   e.Code,
			Detail: fmt.Sprintf("%s: %s", message, e.Detail),
			Status: e.Status,
		}
	}

	// For generic errors, wrap them
	return New("WRAPPED_ERROR", fmt.Sprintf("%s: %v", message, originalError), 500)
}

// WithDetail creates a new error with additional detail
func WithDetail(originalError error, detail string) error {
	if originalError == nil {
		return New("DETAILED_ERROR", detail, 500)
	}

	if e, ok := originalError.(*Error); ok {
		return &Error{
			Id:     e.Id,
			Code:   e.Code,
			Detail: detail,
			Status: e.Status,
		}
	}

	return New("DETAILED_ERROR", detail, 500)
}

// GetHTTPStatus returns the HTTP status code from an error
func GetHTTPStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return int(e.Code)
	}
	return http.StatusInternalServerError
}

// GetErrorID returns the error ID from an error
func GetErrorID(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Id
	}
	return "UNKNOWN_ERROR"
}

// ============================================================================
// Domain-Specific Helper Functions
// ============================================================================

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	return Is(err, ErrCollateralNotFound) || Is(err, NotFound(""))
}

// IsAlreadyExists checks if the error is a "conflict" error
func IsAlreadyExists(err error) bool {
	return Is(err, ErrCollateralExists) || Is(err, Conflict(""))
}

// IsBadRequest checks if the error is a "bad request" error
func IsBadRequest(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code >= 400 && e.Code < 500
	}
	return false
}

// IsServerError checks if the error is a server error (5xx)
func IsServerError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code >= 500 && e.Code < 600
	}
	return false
}

// ============================================================================
// Builder Pattern for Custom Errors
// ============================================================================

// ErrorBuilder provides a fluent interface for building errors
type ErrorBuilder struct {
	id     string
	code   int32
	detail string
}

// NewBuilder creates a new error builder
func NewBuilder() *ErrorBuilder {
	return &ErrorBuilder{
		id:   "CUSTOM_ERROR",
		code: 500,
	}
}

// WithID sets the error ID
func (b *ErrorBuilder) WithID(id string) *ErrorBuilder {
	b.id = id
	return b
}

// WithCode sets the HTTP status code
func (b *ErrorBuilder) WithCode(code int32) *ErrorBuilder {
	b.code = code
	return b
}

// WithDetail sets the error detail message
func (b *ErrorBuilder) WithDetail(detail string) *ErrorBuilder {
	b.detail = detail
	return b
}

// WithDetailf sets the error detail message with formatting
func (b *ErrorBuilder) WithDetailf(format string, args ...interface{}) *ErrorBuilder {
	b.detail = fmt.Sprintf(format, args...)
	return b
}

// Build creates the final error
func (b *ErrorBuilder) Build() error {
	return New(b.id, b.detail, b.code)
}
