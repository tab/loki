// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/repositories/scope.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/repositories/scope.go -destination=internal/app/repositories/scope_mock.go -package=repositories
//

// Package repositories is a generated GoMock package.
package repositories

import (
	context "context"
	models "loki/internal/app/models"
	db "loki/internal/app/repositories/db"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockScopeRepository is a mock of ScopeRepository interface.
type MockScopeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockScopeRepositoryMockRecorder
	isgomock struct{}
}

// MockScopeRepositoryMockRecorder is the mock recorder for MockScopeRepository.
type MockScopeRepositoryMockRecorder struct {
	mock *MockScopeRepository
}

// NewMockScopeRepository creates a new mock instance.
func NewMockScopeRepository(ctrl *gomock.Controller) *MockScopeRepository {
	mock := &MockScopeRepository{ctrl: ctrl}
	mock.recorder = &MockScopeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScopeRepository) EXPECT() *MockScopeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockScopeRepository) Create(ctx context.Context, params db.CreateScopeParams) (*models.Scope, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*models.Scope)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockScopeRepositoryMockRecorder) Create(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScopeRepository)(nil).Create), ctx, params)
}

// CreateUserScope mocks base method.
func (m *MockScopeRepository) CreateUserScope(ctx context.Context, params db.CreateUserScopeParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserScope", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserScope indicates an expected call of CreateUserScope.
func (mr *MockScopeRepositoryMockRecorder) CreateUserScope(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserScope", reflect.TypeOf((*MockScopeRepository)(nil).CreateUserScope), ctx, params)
}

// Delete mocks base method.
func (m *MockScopeRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockScopeRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockScopeRepository)(nil).Delete), ctx, id)
}

// FindById mocks base method.
func (m *MockScopeRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Scope, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, id)
	ret0, _ := ret[0].(*models.Scope)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockScopeRepositoryMockRecorder) FindById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockScopeRepository)(nil).FindById), ctx, id)
}

// FindByName mocks base method.
func (m *MockScopeRepository) FindByName(ctx context.Context, name string) (*models.Scope, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", ctx, name)
	ret0, _ := ret[0].(*models.Scope)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName.
func (mr *MockScopeRepositoryMockRecorder) FindByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockScopeRepository)(nil).FindByName), ctx, name)
}

// FindByUserId mocks base method.
func (m *MockScopeRepository) FindByUserId(ctx context.Context, id uuid.UUID) ([]models.Scope, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserId", ctx, id)
	ret0, _ := ret[0].([]models.Scope)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserId indicates an expected call of FindByUserId.
func (mr *MockScopeRepositoryMockRecorder) FindByUserId(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserId", reflect.TypeOf((*MockScopeRepository)(nil).FindByUserId), ctx, id)
}

// List mocks base method.
func (m *MockScopeRepository) List(ctx context.Context, limit, offset int32) ([]models.Scope, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset)
	ret0, _ := ret[0].([]models.Scope)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockScopeRepositoryMockRecorder) List(ctx, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockScopeRepository)(nil).List), ctx, limit, offset)
}

// Update mocks base method.
func (m *MockScopeRepository) Update(ctx context.Context, params db.UpdateScopeParams) (*models.Scope, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, params)
	ret0, _ := ret[0].(*models.Scope)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockScopeRepositoryMockRecorder) Update(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockScopeRepository)(nil).Update), ctx, params)
}
