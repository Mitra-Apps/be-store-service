// Code generated by MockGen. DO NOT EDIT.
// Source: domain/store/repository/repository.go
//
// Generated by this command:
//
//	mockgen -source=domain/store/repository/repository.go -destination=domain/store/repository/mocks/repository.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/Mitra-Apps/be-store-service/domain/store/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockStoreServiceRepository is a mock of StoreServiceRepository interface.
type MockStoreServiceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStoreServiceRepositoryMockRecorder
}

// MockStoreServiceRepositoryMockRecorder is the mock recorder for MockStoreServiceRepository.
type MockStoreServiceRepositoryMockRecorder struct {
	mock *MockStoreServiceRepository
}

// NewMockStoreServiceRepository creates a new mock instance.
func NewMockStoreServiceRepository(ctrl *gomock.Controller) *MockStoreServiceRepository {
	mock := &MockStoreServiceRepository{ctrl: ctrl}
	mock.recorder = &MockStoreServiceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreServiceRepository) EXPECT() *MockStoreServiceRepositoryMockRecorder {
	return m.recorder
}

// CreateStore mocks base method.
func (m *MockStoreServiceRepository) CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStore", ctx, store)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStore indicates an expected call of CreateStore.
func (mr *MockStoreServiceRepositoryMockRecorder) CreateStore(ctx, store any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStore", reflect.TypeOf((*MockStoreServiceRepository)(nil).CreateStore), ctx, store)
}

// DeleteStores mocks base method.
func (m *MockStoreServiceRepository) DeleteStores(ctx context.Context, storeIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStores", ctx, storeIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStores indicates an expected call of DeleteStores.
func (mr *MockStoreServiceRepositoryMockRecorder) DeleteStores(ctx, storeIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStores", reflect.TypeOf((*MockStoreServiceRepository)(nil).DeleteStores), ctx, storeIDs)
}

// GetStore mocks base method.
func (m *MockStoreServiceRepository) GetStore(ctx context.Context, storeID string) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStore", ctx, storeID)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStore indicates an expected call of GetStore.
func (mr *MockStoreServiceRepositoryMockRecorder) GetStore(ctx, storeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStore", reflect.TypeOf((*MockStoreServiceRepository)(nil).GetStore), ctx, storeID)
}

// GetStoreByUserID mocks base method.
func (m *MockStoreServiceRepository) GetStoreByUserID(ctx context.Context, userID uuid.UUID) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStoreByUserID", ctx, userID)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStoreByUserID indicates an expected call of GetStoreByUserID.
func (mr *MockStoreServiceRepositoryMockRecorder) GetStoreByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStoreByUserID", reflect.TypeOf((*MockStoreServiceRepository)(nil).GetStoreByUserID), ctx, userID)
}

// ListStores mocks base method.
func (m *MockStoreServiceRepository) ListStores(ctx context.Context, page, pageSize int) ([]*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStores", ctx, page, pageSize)
	ret0, _ := ret[0].([]*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStores indicates an expected call of ListStores.
func (mr *MockStoreServiceRepositoryMockRecorder) ListStores(ctx, page, pageSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStores", reflect.TypeOf((*MockStoreServiceRepository)(nil).ListStores), ctx, page, pageSize)
}

// OpenCloseStore mocks base method.
func (m *MockStoreServiceRepository) OpenCloseStore(ctx context.Context, storeID uuid.UUID, isActive bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenCloseStore", ctx, storeID, isActive)
	ret0, _ := ret[0].(error)
	return ret0
}

// OpenCloseStore indicates an expected call of OpenCloseStore.
func (mr *MockStoreServiceRepositoryMockRecorder) OpenCloseStore(ctx, storeID, isActive any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenCloseStore", reflect.TypeOf((*MockStoreServiceRepository)(nil).OpenCloseStore), ctx, storeID, isActive)
}

// UpdateStore mocks base method.
func (m *MockStoreServiceRepository) UpdateStore(ctx context.Context, update *entity.Store) (*entity.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStore", ctx, update)
	ret0, _ := ret[0].(*entity.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStore indicates an expected call of UpdateStore.
func (mr *MockStoreServiceRepositoryMockRecorder) UpdateStore(ctx, update any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStore", reflect.TypeOf((*MockStoreServiceRepository)(nil).UpdateStore), ctx, update)
}

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// UploadImage mocks base method.
func (m *MockStorage) UploadImage(ctx context.Context, image, userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", ctx, image, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockStorageMockRecorder) UploadImage(ctx, image, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockStorage)(nil).UploadImage), ctx, image, userID)
}
