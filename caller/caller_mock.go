// Code generated by MockGen. DO NOT EDIT.
// Source: caller/caller.go

// Package caller is a generated GoMock package.
package caller

import (
	context "context"
	json "encoding/json"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCaller is a mock of Caller interface.
type MockCaller struct {
	ctrl     *gomock.Controller
	recorder *MockCallerMockRecorder
}

// MockCallerMockRecorder is the mock recorder for MockCaller.
type MockCallerMockRecorder struct {
	mock *MockCaller
}

// NewMockCaller creates a new mock instance.
func NewMockCaller(ctrl *gomock.Controller) *MockCaller {
	mock := &MockCaller{ctrl: ctrl}
	mock.recorder = &MockCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCaller) EXPECT() *MockCallerMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockCaller) Call(ctx context.Context, api, method string, args []interface{}, reply interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, api, method, args, reply)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockCallerMockRecorder) Call(ctx, api, method, args, reply interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockCaller)(nil).Call), ctx, api, method, args, reply)
}

// SetCallback mocks base method.
func (m *MockCaller) SetCallback(api, method string, callback func(json.RawMessage)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCallback", api, method, callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCallback indicates an expected call of SetCallback.
func (mr *MockCallerMockRecorder) SetCallback(api, method, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCallback", reflect.TypeOf((*MockCaller)(nil).SetCallback), api, method, callback)
}

// MockCallCloser is a mock of CallCloser interface.
type MockCallCloser struct {
	ctrl     *gomock.Controller
	recorder *MockCallCloserMockRecorder
}

// MockCallCloserMockRecorder is the mock recorder for MockCallCloser.
type MockCallCloserMockRecorder struct {
	mock *MockCallCloser
}

// NewMockCallCloser creates a new mock instance.
func NewMockCallCloser(ctrl *gomock.Controller) *MockCallCloser {
	mock := &MockCallCloser{ctrl: ctrl}
	mock.recorder = &MockCallCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCallCloser) EXPECT() *MockCallCloserMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockCallCloser) Call(ctx context.Context, api, method string, args []interface{}, reply interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, api, method, args, reply)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockCallCloserMockRecorder) Call(ctx, api, method, args, reply interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockCallCloser)(nil).Call), ctx, api, method, args, reply)
}

// Close mocks base method.
func (m *MockCallCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockCallCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCallCloser)(nil).Close))
}

// SetCallback mocks base method.
func (m *MockCallCloser) SetCallback(api, method string, callback func(json.RawMessage)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCallback", api, method, callback)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCallback indicates an expected call of SetCallback.
func (mr *MockCallCloserMockRecorder) SetCallback(api, method, callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCallback", reflect.TypeOf((*MockCallCloser)(nil).SetCallback), api, method, callback)
}
