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

	uuid "github.com/google/uuid"
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

// CreateUser mocks base method.
func (m *MockDatabase) CreateUser(ctx context.Context, params db.CreateUserTokensParams) (*models.User, error) {
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

// CreateUserTokens mocks base method.
func (m *MockDatabase) CreateUserTokens(ctx context.Context, params dto.CreateUserTokensParams) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTokens", ctx, params)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserTokens indicates an expected call of CreateUserTokens.
func (mr *MockDatabaseMockRecorder) CreateUserTokens(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTokens", reflect.TypeOf((*MockDatabase)(nil).CreateUserTokens), ctx, params)
}

// FindUserById mocks base method.
func (m *MockDatabase) FindUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserById", ctx, id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserById indicates an expected call of FindUserById.
func (mr *MockDatabaseMockRecorder) FindUserById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserById", reflect.TypeOf((*MockDatabase)(nil).FindUserById), ctx, id)
}

// FindUserByIdentityNumber mocks base method.
func (m *MockDatabase) FindUserByIdentityNumber(ctx context.Context, identityNumber string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByIdentityNumber", ctx, identityNumber)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByIdentityNumber indicates an expected call of FindUserByIdentityNumber.
func (mr *MockDatabaseMockRecorder) FindUserByIdentityNumber(ctx, identityNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByIdentityNumber", reflect.TypeOf((*MockDatabase)(nil).FindUserByIdentityNumber), ctx, identityNumber)
}
