// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/golibs/params (interfaces: Env)

// Package mock_params is a generated GoMock package.
package mock_params

import (
	gomock "github.com/golang/mock/gomock"
	logging "github.com/qdm12/golibs/logging"
	params "github.com/qdm12/golibs/params"
	url "net/url"
	reflect "reflect"
	time "time"
)

// MockEnv is a mock of Env interface
type MockEnv struct {
	ctrl     *gomock.Controller
	recorder *MockEnvMockRecorder
}

// MockEnvMockRecorder is the mock recorder for MockEnv
type MockEnvMockRecorder struct {
	mock *MockEnv
}

// NewMockEnv creates a new mock instance
func NewMockEnv(ctrl *gomock.Controller) *MockEnv {
	mock := &MockEnv{ctrl: ctrl}
	mock.recorder = &MockEnvMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEnv) EXPECT() *MockEnvMockRecorder {
	return m.recorder
}

// CSVInside mocks base method
func (m *MockEnv) CSVInside(arg0 string, arg1 []string, arg2 ...params.OptionSetter) ([]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CSVInside", varargs...)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CSVInside indicates an expected call of CSVInside
func (mr *MockEnvMockRecorder) CSVInside(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CSVInside", reflect.TypeOf((*MockEnv)(nil).CSVInside), varargs...)
}

// DatabaseDetails mocks base method
func (m *MockEnv) DatabaseDetails() (string, string, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatabaseDetails")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(string)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// DatabaseDetails indicates an expected call of DatabaseDetails
func (mr *MockEnvMockRecorder) DatabaseDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatabaseDetails", reflect.TypeOf((*MockEnv)(nil).DatabaseDetails))
}

// Duration mocks base method
func (m *MockEnv) Duration(arg0 string, arg1 ...params.OptionSetter) (time.Duration, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Duration", varargs...)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Duration indicates an expected call of Duration
func (mr *MockEnvMockRecorder) Duration(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Duration", reflect.TypeOf((*MockEnv)(nil).Duration), varargs...)
}

// ExeDir mocks base method
func (m *MockEnv) ExeDir() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExeDir")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExeDir indicates an expected call of ExeDir
func (mr *MockEnvMockRecorder) ExeDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExeDir", reflect.TypeOf((*MockEnv)(nil).ExeDir))
}

// Get mocks base method
func (m *MockEnv) Get(arg0 string, arg1 ...params.OptionSetter) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockEnvMockRecorder) Get(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockEnv)(nil).Get), varargs...)
}

// GotifyToken mocks base method
func (m *MockEnv) GotifyToken(arg0 ...params.OptionSetter) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GotifyToken", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GotifyToken indicates an expected call of GotifyToken
func (mr *MockEnvMockRecorder) GotifyToken(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GotifyToken", reflect.TypeOf((*MockEnv)(nil).GotifyToken), arg0...)
}

// GotifyURL mocks base method
func (m *MockEnv) GotifyURL(arg0 ...params.OptionSetter) (*url.URL, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GotifyURL", varargs...)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GotifyURL indicates an expected call of GotifyURL
func (mr *MockEnvMockRecorder) GotifyURL(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GotifyURL", reflect.TypeOf((*MockEnv)(nil).GotifyURL), arg0...)
}

// GroupID mocks base method
func (m *MockEnv) GroupID() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GroupID")
	ret0, _ := ret[0].(int)
	return ret0
}

// GroupID indicates an expected call of GroupID
func (mr *MockEnvMockRecorder) GroupID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GroupID", reflect.TypeOf((*MockEnv)(nil).GroupID))
}

// HTTPTimeout mocks base method
func (m *MockEnv) HTTPTimeout(arg0 ...params.OptionSetter) (time.Duration, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HTTPTimeout", varargs...)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HTTPTimeout indicates an expected call of HTTPTimeout
func (mr *MockEnvMockRecorder) HTTPTimeout(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPTimeout", reflect.TypeOf((*MockEnv)(nil).HTTPTimeout), arg0...)
}

// Inside mocks base method
func (m *MockEnv) Inside(arg0 string, arg1 []string, arg2 ...params.OptionSetter) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Inside", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Inside indicates an expected call of Inside
func (mr *MockEnvMockRecorder) Inside(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inside", reflect.TypeOf((*MockEnv)(nil).Inside), varargs...)
}

// Int mocks base method
func (m *MockEnv) Int(arg0 string, arg1 ...params.OptionSetter) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Int", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Int indicates an expected call of Int
func (mr *MockEnvMockRecorder) Int(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int", reflect.TypeOf((*MockEnv)(nil).Int), varargs...)
}

// IntRange mocks base method
func (m *MockEnv) IntRange(arg0 string, arg1, arg2 int, arg3 ...params.OptionSetter) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IntRange", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IntRange indicates an expected call of IntRange
func (mr *MockEnvMockRecorder) IntRange(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IntRange", reflect.TypeOf((*MockEnv)(nil).IntRange), varargs...)
}

// ListeningPort mocks base method
func (m *MockEnv) ListeningPort(arg0 string, arg1 ...params.OptionSetter) (uint16, string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListeningPort", varargs...)
	ret0, _ := ret[0].(uint16)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListeningPort indicates an expected call of ListeningPort
func (mr *MockEnvMockRecorder) ListeningPort(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListeningPort", reflect.TypeOf((*MockEnv)(nil).ListeningPort), varargs...)
}

// LoggerConfig mocks base method
func (m *MockEnv) LoggerConfig() (logging.Encoding, logging.Level, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoggerConfig")
	ret0, _ := ret[0].(logging.Encoding)
	ret1, _ := ret[1].(logging.Level)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LoggerConfig indicates an expected call of LoggerConfig
func (mr *MockEnvMockRecorder) LoggerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoggerConfig", reflect.TypeOf((*MockEnv)(nil).LoggerConfig))
}

// LoggerEncoding mocks base method
func (m *MockEnv) LoggerEncoding(arg0 ...params.OptionSetter) (logging.Encoding, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoggerEncoding", varargs...)
	ret0, _ := ret[0].(logging.Encoding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoggerEncoding indicates an expected call of LoggerEncoding
func (mr *MockEnvMockRecorder) LoggerEncoding(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoggerEncoding", reflect.TypeOf((*MockEnv)(nil).LoggerEncoding), arg0...)
}

// LoggerLevel mocks base method
func (m *MockEnv) LoggerLevel(arg0 ...params.OptionSetter) (logging.Level, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoggerLevel", varargs...)
	ret0, _ := ret[0].(logging.Level)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoggerLevel indicates an expected call of LoggerLevel
func (mr *MockEnvMockRecorder) LoggerLevel(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoggerLevel", reflect.TypeOf((*MockEnv)(nil).LoggerLevel), arg0...)
}

// OnOff mocks base method
func (m *MockEnv) OnOff(arg0 string, arg1 ...params.OptionSetter) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OnOff", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OnOff indicates an expected call of OnOff
func (mr *MockEnvMockRecorder) OnOff(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnOff", reflect.TypeOf((*MockEnv)(nil).OnOff), varargs...)
}

// Path mocks base method
func (m *MockEnv) Path(arg0 string, arg1 ...params.OptionSetter) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Path", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Path indicates an expected call of Path
func (mr *MockEnvMockRecorder) Path(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockEnv)(nil).Path), varargs...)
}

// Port mocks base method
func (m *MockEnv) Port(arg0 string, arg1 ...params.OptionSetter) (uint16, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Port", varargs...)
	ret0, _ := ret[0].(uint16)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Port indicates an expected call of Port
func (mr *MockEnvMockRecorder) Port(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Port", reflect.TypeOf((*MockEnv)(nil).Port), varargs...)
}

// RedisDetails mocks base method
func (m *MockEnv) RedisDetails() (string, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedisDetails")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// RedisDetails indicates an expected call of RedisDetails
func (mr *MockEnvMockRecorder) RedisDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedisDetails", reflect.TypeOf((*MockEnv)(nil).RedisDetails))
}

// RootURL mocks base method
func (m *MockEnv) RootURL(arg0 ...params.OptionSetter) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RootURL", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RootURL indicates an expected call of RootURL
func (mr *MockEnvMockRecorder) RootURL(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RootURL", reflect.TypeOf((*MockEnv)(nil).RootURL), arg0...)
}

// URL mocks base method
func (m *MockEnv) URL(arg0 string, arg1 ...params.OptionSetter) (*url.URL, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "URL", varargs...)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// URL indicates an expected call of URL
func (mr *MockEnvMockRecorder) URL(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockEnv)(nil).URL), varargs...)
}

// UserID mocks base method
func (m *MockEnv) UserID() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserID")
	ret0, _ := ret[0].(int)
	return ret0
}

// UserID indicates an expected call of UserID
func (mr *MockEnvMockRecorder) UserID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserID", reflect.TypeOf((*MockEnv)(nil).UserID))
}

// YesNo mocks base method
func (m *MockEnv) YesNo(arg0 string, arg1 ...params.OptionSetter) (bool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "YesNo", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// YesNo indicates an expected call of YesNo
func (mr *MockEnvMockRecorder) YesNo(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "YesNo", reflect.TypeOf((*MockEnv)(nil).YesNo), varargs...)
}
