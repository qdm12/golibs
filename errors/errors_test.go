package errors

import (
	liberrors "errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPStatus(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		err    error
		status int
	}{
		"nil":               {nil, http.StatusInternalServerError},
		"fmt error":         {fmt.Errorf(""), http.StatusInternalServerError},
		"errors error":      {liberrors.New(""), http.StatusInternalServerError},
		"internal error":    {NewInternal(""), http.StatusInternalServerError},
		"bad requets error": {NewBadRequest(""), http.StatusBadRequest},
		"not found error":   {NewNotFound(""), http.StatusNotFound},
		"conflict error":    {NewConflict(""), http.StatusConflict},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			status := HTTPStatus(tc.err)
			assert.Equal(t, tc.status, status)
		})
	}
}

func Test_NewInternal(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args         []interface{}
		errString    string
		expectedType interface{}
	}{
		"no args":          {nil, "", &InternalError{}},
		"one string":       {[]interface{}{"a"}, "a", &InternalError{}},
		"one fmt error":    {[]interface{}{fmt.Errorf("a")}, "a", &InternalError{}},
		"one errors error": {[]interface{}{liberrors.New("a")}, "a", &InternalError{}},
		"one custom error": {[]interface{}{NewBadRequest("a")}, "a", &BadRequestError{}},
		"two args":         {[]interface{}{"a", 2}, "a 2", &InternalError{}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var err error
			err = NewInternal(tc.args...)
			assert.IsType(t, tc.expectedType, err)
			assert.Equal(t, tc.errString, err.Error())
		})
	}
}

func Test_NewBadRequest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args         []interface{}
		errString    string
		expectedType interface{}
	}{
		"no args":          {nil, "", &BadRequestError{}},
		"one string":       {[]interface{}{"a"}, "a", &BadRequestError{}},
		"one fmt error":    {[]interface{}{fmt.Errorf("a")}, "a", &BadRequestError{}},
		"one errors error": {[]interface{}{liberrors.New("a")}, "a", &BadRequestError{}},
		"one custom error": {[]interface{}{NewInternal("a")}, "a", &InternalError{}},
		"two args":         {[]interface{}{"a", 2}, "a 2", &BadRequestError{}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var err error
			err = NewBadRequest(tc.args...)
			assert.IsType(t, tc.expectedType, err)
			assert.Equal(t, tc.errString, err.Error())
		})
	}
}

func Test_NewNotFound(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args         []interface{}
		errString    string
		expectedType interface{}
	}{
		"no args":          {nil, "", &NotFoundError{}},
		"one string":       {[]interface{}{"a"}, "a", &NotFoundError{}},
		"one fmt error":    {[]interface{}{fmt.Errorf("a")}, "a", &NotFoundError{}},
		"one errors error": {[]interface{}{liberrors.New("a")}, "a", &NotFoundError{}},
		"one custom error": {[]interface{}{NewInternal("a")}, "a", &InternalError{}},
		"two args":         {[]interface{}{"a", 2}, "a 2", &NotFoundError{}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var err error
			err = NewNotFound(tc.args...)
			assert.IsType(t, tc.expectedType, err)
			assert.Equal(t, tc.errString, err.Error())
		})
	}
}

func Test_NewConflict(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args         []interface{}
		errString    string
		expectedType interface{}
	}{
		"no args":          {nil, "", &ConflictError{}},
		"one string":       {[]interface{}{"a"}, "a", &ConflictError{}},
		"one fmt error":    {[]interface{}{fmt.Errorf("a")}, "a", &ConflictError{}},
		"one errors error": {[]interface{}{liberrors.New("a")}, "a", &ConflictError{}},
		"one custom error": {[]interface{}{NewInternal("a")}, "a", &InternalError{}},
		"two args":         {[]interface{}{"a", 2}, "a 2", &ConflictError{}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var err error
			err = NewConflict(tc.args...)
			assert.IsType(t, tc.expectedType, err)
			assert.Equal(t, tc.errString, err.Error())
		})
	}
}
