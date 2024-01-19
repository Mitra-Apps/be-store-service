// Code generated by MockGen. DO NOT EDIT.
// Source: service/service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/Mitra-Apps/be-store-service/domain/store/entity"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateStore mocks base method.
func (m *MockService) CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStore", ctx, store)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStore indicates an expected call of CreateStore.
func (mr *MockServiceMockRecorder) CreateStore(ctx, store interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStore", reflect.TypeOf((*MockService)(nil).CreateStore), ctx, store)
}

// GetStore mocks base method.
func (m *MockService) GetStore(ctx context.Context, storeID string) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStore", ctx, storeID)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStore indicates an expected call of GetStore.
func (mr *MockServiceMockRecorder) GetStore(ctx, storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStore", reflect.TypeOf((*MockService)(nil).GetStore), ctx, storeID)
}

// ListStores mocks base method.
func (m *MockService) ListStores(ctx context.Context) ([]*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStores", ctx)
	ret0, _ := ret[0].([]*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStores indicates an expected call of ListStores.
func (mr *MockServiceMockRecorder) ListStores(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStores", reflect.TypeOf((*MockService)(nil).ListStores), ctx)
}

// OpenCloseStore mocks base method.
func (m *MockService) OpenCloseStore(ctx context.Context, userID uuid.UUID, roleNames []string, storeID string, isActive bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenCloseStore", ctx, userID, roleNames, storeID, isActive)
	ret0, _ := ret[0].(error)
	return ret0
}

// OpenCloseStore indicates an expected call of OpenCloseStore.
func (mr *MockServiceMockRecorder) OpenCloseStore(ctx, userID, roleNames, storeID, isActive interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenCloseStore", reflect.TypeOf((*MockService)(nil).OpenCloseStore), ctx, userID, roleNames, storeID, isActive)
}
