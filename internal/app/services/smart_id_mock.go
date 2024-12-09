// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/services/smart_id.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/services/smart_id.go -destination=internal/app/services/smart_id_mock.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	dto "loki/internal/app/models/dto"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSmartIdProvider is a mock of SmartIdProvider interface.
type MockSmartIdProvider struct {
	ctrl     *gomock.Controller
	recorder *MockSmartIdProviderMockRecorder
	isgomock struct{}
}

// MockSmartIdProviderMockRecorder is the mock recorder for MockSmartIdProvider.
type MockSmartIdProviderMockRecorder struct {
	mock *MockSmartIdProvider
}

// NewMockSmartIdProvider creates a new mock instance.
func NewMockSmartIdProvider(ctrl *gomock.Controller) *MockSmartIdProvider {
	mock := &MockSmartIdProvider{ctrl: ctrl}
	mock.recorder = &MockSmartIdProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSmartIdProvider) EXPECT() *MockSmartIdProviderMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSmartIdProvider) CreateSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*dto.SmartIdProviderSessionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, params)
	ret0, _ := ret[0].(*dto.SmartIdProviderSessionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSmartIdProviderMockRecorder) CreateSession(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSmartIdProvider)(nil).CreateSession), ctx, params)
}

// GetSessionStatus mocks base method.
func (m *MockSmartIdProvider) GetSessionStatus(id uuid.UUID) (*dto.SmartIdProviderSessionStatusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionStatus", id)
	ret0, _ := ret[0].(*dto.SmartIdProviderSessionStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionStatus indicates an expected call of GetSessionStatus.
func (mr *MockSmartIdProviderMockRecorder) GetSessionStatus(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionStatus", reflect.TypeOf((*MockSmartIdProvider)(nil).GetSessionStatus), id)
}
