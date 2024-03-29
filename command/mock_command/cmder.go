// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/golibs/command (interfaces: RunStarter)

// Package mock_command is a generated GoMock package.
package mock_command

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	command "github.com/qdm12/golibs/command"
)

// MockRunStarter is a mock of RunStarter interface.
type MockRunStarter struct {
	ctrl     *gomock.Controller
	recorder *MockRunStarterMockRecorder
}

// MockRunStarterMockRecorder is the mock recorder for MockRunStarter.
type MockRunStarterMockRecorder struct {
	mock *MockRunStarter
}

// NewMockRunStarter creates a new mock instance.
func NewMockRunStarter(ctrl *gomock.Controller) *MockRunStarter {
	mock := &MockRunStarter{ctrl: ctrl}
	mock.recorder = &MockRunStarterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRunStarter) EXPECT() *MockRunStarterMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockRunStarter) Run(arg0 command.ExecCmd) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockRunStarterMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockRunStarter)(nil).Run), arg0)
}

// Start mocks base method.
func (m *MockRunStarter) Start(arg0 command.ExecCmd) (chan string, chan string, chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(chan string)
	ret1, _ := ret[1].(chan string)
	ret2, _ := ret[2].(chan error)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Start indicates an expected call of Start.
func (mr *MockRunStarterMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRunStarter)(nil).Start), arg0)
}
