package errors

import (
	"net/http"

	"github.com/qdm12/golibs/format"
)

type InternalError struct{ s string }
type BadRequestError struct{ s string }
type NotFoundError struct{ s string }
type ConflictError struct{ s string }

func (e *InternalError) Error() string   { return e.s }
func (e *BadRequestError) Error() string { return e.s }
func (e *NotFoundError) Error() string   { return e.s }
func (e *ConflictError) Error() string   { return e.s }

func NewInternal(args ...interface{}) error {
	if len(args) == 1 && isErrorCustom(args[0]) {
		return args[0].(error)
	}
	return &InternalError{format.ArgsToString(args...)}
}

func NewBadRequest(args ...interface{}) error {
	if len(args) == 1 && isErrorCustom(args[0]) {
		return args[0].(error)
	}
	return &BadRequestError{format.ArgsToString(args...)}
}

func NewNotFound(args ...interface{}) error {
	if len(args) == 1 && isErrorCustom(args[0]) {
		return args[0].(error)
	}
	return &NotFoundError{format.ArgsToString(args...)}
}

func NewConflict(args ...interface{}) error {
	if len(args) == 1 && isErrorCustom(args[0]) {
		return args[0].(error)
	}
	return &ConflictError{format.ArgsToString(args...)}
}

func isErrorCustom(err interface{}) bool {
	switch err.(type) {
	case *InternalError, *BadRequestError, *NotFoundError, *ConflictError:
		return true
	}
	return false
}

func HTTPStatus(err error) int {
	switch err.(type) {
	case *InternalError:
		return http.StatusInternalServerError
	case *BadRequestError:
		return http.StatusBadRequest
	case *NotFoundError:
		return http.StatusNotFound
	case *ConflictError:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
