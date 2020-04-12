// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/golibs/crypto (interfaces: Crypto)

// Package mock_crypto is a generated GoMock package.
package mock_crypto

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCrypto is a mock of Crypto interface
type MockCrypto struct {
	ctrl     *gomock.Controller
	recorder *MockCryptoMockRecorder
}

// MockCryptoMockRecorder is the mock recorder for MockCrypto
type MockCryptoMockRecorder struct {
	mock *MockCrypto
}

// NewMockCrypto creates a new mock instance
func NewMockCrypto(ctrl *gomock.Controller) *MockCrypto {
	mock := &MockCrypto{ctrl: ctrl}
	mock.recorder = &MockCryptoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCrypto) EXPECT() *MockCryptoMockRecorder {
	return m.recorder
}

// Argon2ID mocks base method
func (m *MockCrypto) Argon2ID(arg0 []byte, arg1, arg2 uint32) [64]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Argon2ID", arg0, arg1, arg2)
	ret0, _ := ret[0].([64]byte)
	return ret0
}

// Argon2ID indicates an expected call of Argon2ID
func (mr *MockCryptoMockRecorder) Argon2ID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Argon2ID", reflect.TypeOf((*MockCrypto)(nil).Argon2ID), arg0, arg1, arg2)
}

// Checksumize mocks base method
func (m *MockCrypto) Checksumize(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checksumize", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Checksumize indicates an expected call of Checksumize
func (mr *MockCryptoMockRecorder) Checksumize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checksumize", reflect.TypeOf((*MockCrypto)(nil).Checksumize), arg0)
}

// Dechecksumize mocks base method
func (m *MockCrypto) Dechecksumize(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dechecksumize", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Dechecksumize indicates an expected call of Dechecksumize
func (mr *MockCryptoMockRecorder) Dechecksumize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dechecksumize", reflect.TypeOf((*MockCrypto)(nil).Dechecksumize), arg0)
}

// DecryptAES256 mocks base method
func (m *MockCrypto) DecryptAES256(arg0 []byte, arg1 [32]byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecryptAES256", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecryptAES256 indicates an expected call of DecryptAES256
func (mr *MockCryptoMockRecorder) DecryptAES256(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecryptAES256", reflect.TypeOf((*MockCrypto)(nil).DecryptAES256), arg0, arg1)
}

// EncryptAES256 mocks base method
func (m *MockCrypto) EncryptAES256(arg0 []byte, arg1 [32]byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptAES256", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptAES256 indicates an expected call of EncryptAES256
func (mr *MockCryptoMockRecorder) EncryptAES256(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptAES256", reflect.TypeOf((*MockCrypto)(nil).EncryptAES256), arg0, arg1)
}

// NewSalt mocks base method
func (m *MockCrypto) NewSalt() ([32]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSalt")
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSalt indicates an expected call of NewSalt
func (mr *MockCryptoMockRecorder) NewSalt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSalt", reflect.TypeOf((*MockCrypto)(nil).NewSalt))
}

// ShakeSum256 mocks base method
func (m *MockCrypto) ShakeSum256(arg0 []byte) ([64]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShakeSum256", arg0)
	ret0, _ := ret[0].([64]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShakeSum256 indicates an expected call of ShakeSum256
func (mr *MockCryptoMockRecorder) ShakeSum256(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShakeSum256", reflect.TypeOf((*MockCrypto)(nil).ShakeSum256), arg0)
}

// SignEd25519 mocks base method
func (m *MockCrypto) SignEd25519(arg0 []byte, arg1 [64]byte) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignEd25519", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// SignEd25519 indicates an expected call of SignEd25519
func (mr *MockCryptoMockRecorder) SignEd25519(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignEd25519", reflect.TypeOf((*MockCrypto)(nil).SignEd25519), arg0, arg1)
}
