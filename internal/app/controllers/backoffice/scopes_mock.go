// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/controllers/backoffice/scopes.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/controllers/backoffice/scopes.go -destination=internal/app/controllers/backoffice/scopes_mock.go -package=backoffice
//

// Package backoffice is a generated GoMock package.
package backoffice

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockScopesController is a mock of ScopesController interface.
type MockScopesController struct {
	ctrl     *gomock.Controller
	recorder *MockScopesControllerMockRecorder
	isgomock struct{}
}

// MockScopesControllerMockRecorder is the mock recorder for MockScopesController.
type MockScopesControllerMockRecorder struct {
	mock *MockScopesController
}

// NewMockScopesController creates a new mock instance.
func NewMockBackofficeScopesController(ctrl *gomock.Controller) *MockScopesController {
	mock := &MockScopesController{ctrl: ctrl}
	mock.recorder = &MockScopesControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScopesController) EXPECT() *MockScopesControllerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockScopesController) Create(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Create", w, r)
}

// Create indicates an expected call of Create.
func (mr *MockScopesControllerMockRecorder) Create(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScopesController)(nil).Create), w, r)
}

// Delete mocks base method.
func (m *MockScopesController) Delete(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", w, r)
}

// Delete indicates an expected call of Delete.
func (mr *MockScopesControllerMockRecorder) Delete(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockScopesController)(nil).Delete), w, r)
}

// Get mocks base method.
func (m *MockScopesController) Get(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Get", w, r)
}

// Get indicates an expected call of Get.
func (mr *MockScopesControllerMockRecorder) Get(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockScopesController)(nil).Get), w, r)
}

// List mocks base method.
func (m *MockScopesController) List(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "List", w, r)
}

// List indicates an expected call of List.
func (mr *MockScopesControllerMockRecorder) List(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockScopesController)(nil).List), w, r)
}

// Update mocks base method.
func (m *MockScopesController) Update(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", w, r)
}

// Update indicates an expected call of Update.
func (mr *MockScopesControllerMockRecorder) Update(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockScopesController)(nil).Update), w, r)
}
