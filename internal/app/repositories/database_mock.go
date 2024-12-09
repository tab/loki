// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/repositories/database.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/repositories/database.go -destination=internal/app/repositories/database_mock.go -package=repositories
//

// Package repositories is a generated GoMock package.
package repositories

import (
	context "context"
	models "loki/internal/app/models"
	dto "loki/internal/app/models/dto"
	db "loki/internal/app/repositories/db"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
	isgomock struct{}
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// CreateOrUpdateUserWithTokens mocks base method.
func (m *MockDatabase) CreateOrUpdateUserWithTokens(ctx context.Context, params dto.CreateUserParams) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateUserWithTokens", ctx, params)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrUpdateUserWithTokens indicates an expected call of CreateOrUpdateUserWithTokens.
func (mr *MockDatabaseMockRecorder) CreateOrUpdateUserWithTokens(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateUserWithTokens", reflect.TypeOf((*MockDatabase)(nil).CreateOrUpdateUserWithTokens), ctx, params)
}

// CreateUser mocks base method.
func (m *MockDatabase) CreateUser(ctx context.Context, params db.CreateUserParams) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, params)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockDatabaseMockRecorder) CreateUser(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockDatabase)(nil).CreateUser), ctx, params)
}