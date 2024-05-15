// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	base_model "github.com/Mitra-Apps/be-store-service/domain/base_model"
	entity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	entity0 "github.com/Mitra-Apps/be-store-service/domain/store/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
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
func (m *MockService) CreateStore(ctx context.Context, store *entity0.Store) (*entity0.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStore", ctx, store)
	ret0, _ := ret[0].(*entity0.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStore indicates an expected call of CreateStore.
func (mr *MockServiceMockRecorder) CreateStore(ctx, store interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStore", reflect.TypeOf((*MockService)(nil).CreateStore), ctx, store)
}

// DeleteProductById mocks base method.
func (m *MockService) DeleteProductById(ctx context.Context, userId, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProductById", ctx, userId, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProductById indicates an expected call of DeleteProductById.
func (mr *MockServiceMockRecorder) DeleteProductById(ctx, userId, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProductById", reflect.TypeOf((*MockService)(nil).DeleteProductById), ctx, userId, id)
}

// DeleteStores mocks base method.
func (m *MockService) DeleteStores(ctx context.Context, storeIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStores", ctx, storeIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStores indicates an expected call of DeleteStores.
func (mr *MockServiceMockRecorder) DeleteStores(ctx, storeIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStores", reflect.TypeOf((*MockService)(nil).DeleteStores), ctx, storeIDs)
}

// GetProductById mocks base method.
func (m *MockService) GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductById", ctx, id)
	ret0, _ := ret[0].(*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductById indicates an expected call of GetProductById.
func (mr *MockServiceMockRecorder) GetProductById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductById", reflect.TypeOf((*MockService)(nil).GetProductById), ctx, id)
}

// GetProductCategories mocks base method.
func (m *MockService) GetProductCategories(ctx context.Context, isIncludeDeactivated bool) ([]*entity.ProductCategory, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductCategories", ctx, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.ProductCategory)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProductCategories indicates an expected call of GetProductCategories.
func (mr *MockServiceMockRecorder) GetProductCategories(ctx, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductCategories", reflect.TypeOf((*MockService)(nil).GetProductCategories), ctx, isIncludeDeactivated)
}

// GetProductTypes mocks base method.
func (m *MockService) GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) ([]*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTypes", ctx, productCategoryID, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTypes indicates an expected call of GetProductTypes.
func (mr *MockServiceMockRecorder) GetProductTypes(ctx, productCategoryID, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTypes", reflect.TypeOf((*MockService)(nil).GetProductTypes), ctx, productCategoryID, isIncludeDeactivated)
}

// GetProductsByStoreId mocks base method.
func (m *MockService) GetProductsByStoreId(ctx context.Context, page, limit int32, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) ([]*entity.Product, base_model.Pagination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByStoreId", ctx, page, limit, storeID, productTypeId, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.Product)
	ret1, _ := ret[1].(base_model.Pagination)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProductsByStoreId indicates an expected call of GetProductsByStoreId.
func (mr *MockServiceMockRecorder) GetProductsByStoreId(ctx, page, limit, storeID, productTypeId, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByStoreId", reflect.TypeOf((*MockService)(nil).GetProductsByStoreId), ctx, page, limit, storeID, productTypeId, isIncludeDeactivated)
}

// GetStore mocks base method.
func (m *MockService) GetStore(ctx context.Context, storeID string) (*entity0.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStore", ctx, storeID)
	ret0, _ := ret[0].(*entity0.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStore indicates an expected call of GetStore.
func (mr *MockServiceMockRecorder) GetStore(ctx, storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStore", reflect.TypeOf((*MockService)(nil).GetStore), ctx, storeID)
}

// GetStoreByUserID mocks base method.
func (m *MockService) GetStoreByUserID(ctx context.Context, userID uuid.UUID) (*entity0.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStoreByUserID", ctx, userID)
	ret0, _ := ret[0].(*entity0.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStoreByUserID indicates an expected call of GetStoreByUserID.
func (mr *MockServiceMockRecorder) GetStoreByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStoreByUserID", reflect.TypeOf((*MockService)(nil).GetStoreByUserID), ctx, userID)
}

// ListStores mocks base method.
func (m *MockService) ListStores(ctx context.Context, page, limit int32) ([]*entity0.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStores", ctx, page, limit)
	ret0, _ := ret[0].([]*entity0.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStores indicates an expected call of ListStores.
func (mr *MockServiceMockRecorder) ListStores(ctx, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStores", reflect.TypeOf((*MockService)(nil).ListStores), ctx, page, limit)
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

// UpdateProductCategory mocks base method.
func (m *MockService) UpdateProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductCategory", ctx, prodCategory)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductCategory indicates an expected call of UpdateProductCategory.
func (mr *MockServiceMockRecorder) UpdateProductCategory(ctx, prodCategory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductCategory", reflect.TypeOf((*MockService)(nil).UpdateProductCategory), ctx, prodCategory)
}

// UpdateStore mocks base method.
func (m *MockService) UpdateStore(ctx context.Context, storeID string, update *entity0.Store) (*entity0.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStore", ctx, storeID, update)
	ret0, _ := ret[0].(*entity0.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStore indicates an expected call of UpdateStore.
func (mr *MockServiceMockRecorder) UpdateStore(ctx, storeID, update interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStore", reflect.TypeOf((*MockService)(nil).UpdateStore), ctx, storeID, update)
}

// UpsertProductCategory mocks base method.
func (m *MockService) UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductCategory", ctx, prodCategory)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductCategory indicates an expected call of UpsertProductCategory.
func (mr *MockServiceMockRecorder) UpsertProductCategory(ctx, prodCategory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProductCategory", reflect.TypeOf((*MockService)(nil).UpsertProductCategory), ctx, prodCategory)
}

// UpsertProductType mocks base method.
func (m *MockService) UpsertProductType(ctx context.Context, prodType *entity.ProductType, isUpdate bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductType", ctx, prodType, isUpdate)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductType indicates an expected call of UpsertProductType.
func (mr *MockServiceMockRecorder) UpsertProductType(ctx, prodType, isUpdate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProductType", reflect.TypeOf((*MockService)(nil).UpsertProductType), ctx, prodType, isUpdate)
}

// UpsertProducts mocks base method.
func (m *MockService) UpsertProducts(ctx context.Context, userID uuid.UUID, roleNames []string, storeID uuid.UUID, isUpdate bool, products ...*entity.Product) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, userID, roleNames, storeID, isUpdate}
	for _, a := range products {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertProducts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProducts indicates an expected call of UpsertProducts.
func (mr *MockServiceMockRecorder) UpsertProducts(ctx, userID, roleNames, storeID, isUpdate interface{}, products ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, userID, roleNames, storeID, isUpdate}, products...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProducts", reflect.TypeOf((*MockService)(nil).UpsertProducts), varargs...)
}
