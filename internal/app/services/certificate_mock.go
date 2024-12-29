// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/services/certificate.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/services/certificate.go -destination=internal/app/services/certificate_mock.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCertificate is a mock of Certificate interface.
type MockCertificate struct {
	ctrl     *gomock.Controller
	recorder *MockCertificateMockRecorder
	isgomock struct{}
}

// MockCertificateMockRecorder is the mock recorder for MockCertificate.
type MockCertificateMockRecorder struct {
	mock *MockCertificate
}

// NewMockCertificate creates a new mock instance.
func NewMockCertificate(ctrl *gomock.Controller) *MockCertificate {
	mock := &MockCertificate{ctrl: ctrl}
	mock.recorder = &MockCertificateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCertificate) EXPECT() *MockCertificateMockRecorder {
	return m.recorder
}

// Extract mocks base method.
func (m *MockCertificate) Extract(value string) (*CertificatePayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", value)
	ret0, _ := ret[0].(*CertificatePayload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Extract indicates an expected call of Extract.
func (mr *MockCertificateMockRecorder) Extract(value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*MockCertificate)(nil).Extract), value)
}
