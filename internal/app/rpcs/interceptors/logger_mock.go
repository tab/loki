// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/rpcs/interceptors/logger.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/rpcs/interceptors/logger.go -destination=internal/app/rpcs/interceptors/logger_mock.go -package=interceptors
//

// Package interceptors is a generated GoMock package.
package interceptors

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockLoggerInterceptor is a mock of LoggerInterceptor interface.
type MockLoggerInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerInterceptorMockRecorder
	isgomock struct{}
}

// MockLoggerInterceptorMockRecorder is the mock recorder for MockLoggerInterceptor.
type MockLoggerInterceptorMockRecorder struct {
	mock *MockLoggerInterceptor
}

// NewMockLoggerInterceptor creates a new mock instance.
func NewMockLoggerInterceptor(ctrl *gomock.Controller) *MockLoggerInterceptor {
	mock := &MockLoggerInterceptor{ctrl: ctrl}
	mock.recorder = &MockLoggerInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerInterceptor) EXPECT() *MockLoggerInterceptorMockRecorder {
	return m.recorder
}

// Log mocks base method.
func (m *MockLoggerInterceptor) Log() grpc.UnaryServerInterceptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(grpc.UnaryServerInterceptor)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockLoggerInterceptorMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockLoggerInterceptor)(nil).Log))
}
