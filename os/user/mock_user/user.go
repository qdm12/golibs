// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/golibs/os/user (interfaces: OSUser)

// Package mock_user is a generated GoMock package.
package mock_user

import (
	gomock "github.com/golang/mock/gomock"
	user "os/user"
	reflect "reflect"
)

// MockOSUser is a mock of OSUser interface
type MockOSUser struct {
	ctrl     *gomock.Controller
	recorder *MockOSUserMockRecorder
}

// MockOSUserMockRecorder is the mock recorder for MockOSUser
type MockOSUserMockRecorder struct {
	mock *MockOSUser
}

// NewMockOSUser creates a new mock instance
func NewMockOSUser(ctrl *gomock.Controller) *MockOSUser {
	mock := &MockOSUser{ctrl: ctrl}
	mock.recorder = &MockOSUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOSUser) EXPECT() *MockOSUserMockRecorder {
	return m.recorder
}

// Lookup mocks base method
func (m *MockOSUser) Lookup(arg0 string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lookup", arg0)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup
func (mr *MockOSUserMockRecorder) Lookup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lookup", reflect.TypeOf((*MockOSUser)(nil).Lookup), arg0)
}

// LookupID mocks base method
func (m *MockOSUser) LookupID(arg0 string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupID", arg0)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupID indicates an expected call of LookupID
func (mr *MockOSUserMockRecorder) LookupID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupID", reflect.TypeOf((*MockOSUser)(nil).LookupID), arg0)
}