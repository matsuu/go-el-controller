// Code generated by MockGen. DO NOT EDIT.
// Source: wisun.go

// Package wisun is a generated GoMock package.
package wisun

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSerialClient is a mock of SerialClient interface
type MockSerialClient struct {
	ctrl     *gomock.Controller
	recorder *MockSerialClientMockRecorder
}

// MockSerialClientMockRecorder is the mock recorder for MockSerialClient
type MockSerialClientMockRecorder struct {
	mock *MockSerialClient
}

// NewMockSerialClient creates a new mock instance
func NewMockSerialClient(ctrl *gomock.Controller) *MockSerialClient {
	mock := &MockSerialClient{ctrl: ctrl}
	mock.recorder = &MockSerialClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSerialClient) EXPECT() *MockSerialClientMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockSerialClient) Send(in []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockSerialClientMockRecorder) Send(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSerialClient)(nil).Send), in)
}

// Recv mocks base method
func (m *MockSerialClient) Recv() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockSerialClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockSerialClient)(nil).Recv))
}

// Close mocks base method
func (m *MockSerialClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockSerialClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSerialClient)(nil).Close))
}
