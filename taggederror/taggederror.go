package taggederror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// RootTaggedErrors (mapping to HTTP status codes)
	ErrNotFound = &TaggedError{
		err:  errors.New("not found"),
		tag:  "NOT_FOUND",
		code: http.StatusNotFound,
	}
	ErrBadRequest = &TaggedError{
		err:  errors.New("bad request"),
		tag:  "BAD_REQUEST",
		code: http.StatusBadRequest,
	}
	ErrUnauthorized = &TaggedError{
		err:  errors.New("unauthorized"),
		tag:  "UNAUTHORIZED",
		code: http.StatusUnauthorized,
	}
	ErrForbidden = &TaggedError{
		err:  errors.New("forbidden"),
		tag:  "FORBIDDEN",
		code: http.StatusForbidden,
	}
	ErrInternal = &TaggedError{
		err:  errors.New("internal server error"),
		tag:  "INTERNAL_SERVER_ERROR",
		code: http.StatusInternalServerError,
	}
)

// A TaggedError is an error implementation that in a nester wrapping provides the tag of the most inner child and the code of the root TaggedError.
// The tag is useful when returning a business related error in an api.
// Example wrapping:
//
//	var ErrProductInUse = taggederror.ErrBadRequest.Wrap(
//		taggederror.New(
//			errors.New("product in use by at least one shop"), // Business Related Error
//			"PRODUCT_IN_USE",
//		),
//	)
//
// returning this error you will get the following result:
//
//	Error(): "bad request: product in use by at least one shop"
//	Tag(): "PRODUCT_IN_USE"
//	Code(): 401
//
// You can also wrap other error types inside. For example:
//
//	var ErrDBNotHandled = taggederror.ErrInternal.Wrap(
//		taggederror.New(
//			errors.New("unhandled database error"),
//			"UNHANDLED_DB_ERROR",
//		),
//	)
//
// Somewhere in DB layer:
//
//	// Wrapping an unexpected error from DB. e.g. unsupported datatype
//	// You can also decide to hide the error and simply return the ErrDBNotHandled
//	return ErrDBNotHandled.Wrap(err)
//
// Final result:
//
//	Error(): "internal server error: unhandled database error: unsupported data type time.Time ..."
//	Tag(): "UNHANDLED_DB_ERROR"
//	Code(): 500
type TaggedError struct {
	err  error
	tag  string
	code int
}

// NewRoot supports adding errors with status codes not available in this package (e.g. 409)
func NewRoot(err error, tag string, code int) *TaggedError {
	return &TaggedError{
		err:  err,
		tag:  tag,
		code: code,
	}
}

// New is meant to either wrapped directly inside a RootTaggedError or indirectly through another TaggedError
func New(err error, tag string) *TaggedError {
	return &TaggedError{
		err: err,
		tag: tag,
	}
}

// Returns the underlying error's Error()
func (e *TaggedError) Error() string {
	return e.err.Error()
}

func (e *TaggedError) Tag() string {
	return e.tag
}

// Status Code of the error. Is zero if the root error is not a RootTaggedError
func (e *TaggedError) Code() int {
	return e.code
}

// Wrap an error with this TaggedError, if the given error is a TaggedError, its tag will be used.
func (e *TaggedError) Wrap(err error) *TaggedError {
	if err, ok := err.(*TaggedError); ok {
		return &TaggedError{
			err:  fmt.Errorf("%w: %w", e.err, err.err),
			tag:  err.tag,
			code: e.code,
		}
	}
	return &TaggedError{
		err:  fmt.Errorf("%w: %w", e.err, err),
		tag:  e.tag,
		code: e.code,
	}
}

// Returns true if the errors.Is() on the underlying error returns true.
// If the given error is a TaggedError, its underlying error will be used.
func (e *TaggedError) Is(target error) bool {
	if target == nil {
		return e == nil
	}
	if target, ok := target.(*TaggedError); ok {
		return errors.Is(e.err, target.err)
	}
	return errors.Is(e.err, target)
}

// For taggederrors Runs err.Is() and returns errors.Is() otherwise.
func Is(err error, target error) bool {
	if err, ok := err.(*TaggedError); ok {
		return err.Is(target)
	}
	return errors.Is(err, target)
}
