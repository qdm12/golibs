// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/golibs/logging (interfaces: Logger)

// Package mock_logging is a generated GoMock package.
package mock_logging

import (
	gomock "github.com/golang/mock/gomock"
	logging "github.com/qdm12/golibs/logging"
	reflect "reflect"
)

// MockLogger is a mock of Logger interface
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method
func (m *MockLogger) Debug(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug
func (mr *MockLoggerMockRecorder) Debug(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), arg0...)
}

// Error mocks base method
func (m *MockLogger) Error(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error
func (mr *MockLoggerMockRecorder) Error(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), arg0...)
}

// Info mocks base method
func (m *MockLogger) Info(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info
func (mr *MockLoggerMockRecorder) Info(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), arg0...)
}

// SetPrefix mocks base method
func (m *MockLogger) SetPrefix(arg0 string) logging.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPrefix", arg0)
	ret0, _ := ret[0].(logging.Logger)
	return ret0
}

// SetPrefix indicates an expected call of SetPrefix
func (mr *MockLoggerMockRecorder) SetPrefix(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPrefix", reflect.TypeOf((*MockLogger)(nil).SetPrefix), arg0)
}

// Sync mocks base method
func (m *MockLogger) Sync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync")
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync
func (mr *MockLoggerMockRecorder) Sync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockLogger)(nil).Sync))
}

// Warn mocks base method
func (m *MockLogger) Warn(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn
func (mr *MockLoggerMockRecorder) Warn(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogger)(nil).Warn), arg0...)
}

// WithPrefix mocks base method
func (m *MockLogger) WithPrefix(arg0 string) logging.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithPrefix", arg0)
	ret0, _ := ret[0].(logging.Logger)
	return ret0
}

// WithPrefix indicates an expected call of WithPrefix
func (mr *MockLoggerMockRecorder) WithPrefix(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithPrefix", reflect.TypeOf((*MockLogger)(nil).WithPrefix), arg0)
}
