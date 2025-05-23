// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/services/sessions.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/services/sessions.go -destination=internal/app/services/sessions_mock.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	models "loki/internal/app/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSessions is a mock of Sessions interface.
type MockSessions struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsMockRecorder
	isgomock struct{}
}

// MockSessionsMockRecorder is the mock recorder for MockSessions.
type MockSessionsMockRecorder struct {
	mock *MockSessions
}

// NewMockSessions creates a new mock instance.
func NewMockSessions(ctrl *gomock.Controller) *MockSessions {
	mock := &MockSessions{ctrl: ctrl}
	mock.recorder = &MockSessionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessions) EXPECT() *MockSessionsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessions) Create(ctx context.Context, params *models.CreateSessionParams) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionsMockRecorder) Create(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessions)(nil).Create), ctx, params)
}

// Delete mocks base method.
func (m *MockSessions) Delete(ctx context.Context, sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionsMockRecorder) Delete(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessions)(nil).Delete), ctx, sessionId)
}

// FindById mocks base method.
func (m *MockSessions) FindById(ctx context.Context, sessionId string) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, sessionId)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockSessionsMockRecorder) FindById(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockSessions)(nil).FindById), ctx, sessionId)
}

// Update mocks base method.
func (m *MockSessions) Update(ctx context.Context, params *models.UpdateSessionParams) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, params)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSessionsMockRecorder) Update(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSessions)(nil).Update), ctx, params)
}
