// Code generated by MockGen. DO NOT EDIT.
// Source: serial.go

// Package wisun is a generated GoMock package.
package wisun

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSerial is a mock of Serial interface
type MockSerial struct {
	ctrl     *gomock.Controller
	recorder *MockSerialMockRecorder
}

// MockSerialMockRecorder is the mock recorder for MockSerial
type MockSerialMockRecorder struct {
	mock *MockSerial
}

// NewMockSerial creates a new mock instance
func NewMockSerial(ctrl *gomock.Controller) *MockSerial {
	mock := &MockSerial{ctrl: ctrl}
	mock.recorder = &MockSerialMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSerial) EXPECT() *MockSerialMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockSerial) Send(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockSerialMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSerial)(nil).Send), arg0)
}

// Recv mocks base method
func (m *MockSerial) Recv() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockSerialMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockSerial)(nil).Recv))
}

// Close mocks base method
func (m *MockSerial) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockSerialMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSerial)(nil).Close))
}
