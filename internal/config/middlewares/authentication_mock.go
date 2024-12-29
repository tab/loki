// Code generated by MockGen. DO NOT EDIT.
// Source: internal/config/middlewares/authentication.go
//
// Generated by this command:
//
//	mockgen -source=internal/config/middlewares/authentication.go -destination=internal/config/middlewares/authentication_mock.go -package=middlewares
//

// Package middlewares is a generated GoMock package.
package middlewares

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthMiddleware is a mock of AuthMiddleware interface.
type MockAuthMiddleware struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMiddlewareMockRecorder
	isgomock struct{}
}

// MockAuthMiddlewareMockRecorder is the mock recorder for MockAuthMiddleware.
type MockAuthMiddlewareMockRecorder struct {
	mock *MockAuthMiddleware
}

// NewMockAuthMiddleware creates a new mock instance.
func NewMockAuthMiddleware(ctrl *gomock.Controller) *MockAuthMiddleware {
	mock := &MockAuthMiddleware{ctrl: ctrl}
	mock.recorder = &MockAuthMiddlewareMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthMiddleware) EXPECT() *MockAuthMiddlewareMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", next)
	ret0, _ := ret[0].(http.Handler)
	return ret0
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockAuthMiddlewareMockRecorder) Authenticate(next any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAuthMiddleware)(nil).Authenticate), next)
}
