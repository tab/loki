// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/workers/smartid.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/workers/smartid.go -destination=internal/app/workers/smartid_mock.go -package=workers
//

// Package workers is a generated GoMock package.
package workers

import (
	context "context"
	models "loki/internal/app/models"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSmartIdWorker is a mock of SmartIdWorker interface.
type MockSmartIdWorker struct {
	ctrl     *gomock.Controller
	recorder *MockSmartIdWorkerMockRecorder
	isgomock struct{}
}

// MockSmartIdWorkerMockRecorder is the mock recorder for MockSmartIdWorker.
type MockSmartIdWorkerMockRecorder struct {
	mock *MockSmartIdWorker
}

// NewMockSmartIdWorker creates a new mock instance.
func NewMockSmartIdWorker(ctrl *gomock.Controller) *MockSmartIdWorker {
	mock := &MockSmartIdWorker{ctrl: ctrl}
	mock.recorder = &MockSmartIdWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSmartIdWorker) EXPECT() *MockSmartIdWorkerMockRecorder {
	return m.recorder
}

// Perform mocks base method.
func (m *MockSmartIdWorker) Perform(ctx context.Context, id uuid.UUID, traceId string) *models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Perform", ctx, id, traceId)
	ret0, _ := ret[0].(*models.Session)
	return ret0
}

// Perform indicates an expected call of Perform.
func (mr *MockSmartIdWorkerMockRecorder) Perform(ctx, id, traceId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Perform", reflect.TypeOf((*MockSmartIdWorker)(nil).Perform), ctx, id, traceId)
}
