// Code generated by MockGen. DO NOT EDIT.
// Source: domain/product/repository/repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// GetProductById mocks base method.
func (m *MockProductRepository) GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductById", ctx, id)
	ret0, _ := ret[0].(*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductById indicates an expected call of GetProductById.
func (mr *MockProductRepositoryMockRecorder) GetProductById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductById", reflect.TypeOf((*MockProductRepository)(nil).GetProductById), ctx, id)
}

// GetProductCategories mocks base method.
func (m *MockProductRepository) GetProductCategories(ctx context.Context, isIncludeDeactivated bool) ([]*entity.ProductCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductCategories", ctx, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.ProductCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductCategories indicates an expected call of GetProductCategories.
func (mr *MockProductRepositoryMockRecorder) GetProductCategories(ctx, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductCategories", reflect.TypeOf((*MockProductRepository)(nil).GetProductCategories), ctx, isIncludeDeactivated)
}

// GetProductTypes mocks base method.
func (m *MockProductRepository) GetProductTypes(ctx context.Context, productCategoryID uuid.UUID, isIncludeDeactivated bool) ([]*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTypes", ctx, productCategoryID, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTypes indicates an expected call of GetProductTypes.
func (mr *MockProductRepositoryMockRecorder) GetProductTypes(ctx, productCategoryID, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTypes", reflect.TypeOf((*MockProductRepository)(nil).GetProductTypes), ctx, productCategoryID, isIncludeDeactivated)
}

// GetProductsByStoreId mocks base method.
func (m *MockProductRepository) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *uuid.UUID, isIncludeDeactivated bool) ([]*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByStoreId", ctx, storeID, productTypeId, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsByStoreId indicates an expected call of GetProductsByStoreId.
func (mr *MockProductRepositoryMockRecorder) GetProductsByStoreId(ctx, storeID, productTypeId, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByStoreId", reflect.TypeOf((*MockProductRepository)(nil).GetProductsByStoreId), ctx, storeID, productTypeId, isIncludeDeactivated)
}

// GetProductsByStoreIdAndNames mocks base method.
func (m *MockProductRepository) GetProductsByStoreIdAndNames(ctx context.Context, storeID uuid.UUID, names []string) ([]*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByStoreIdAndNames", ctx, storeID, names)
	ret0, _ := ret[0].([]*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsByStoreIdAndNames indicates an expected call of GetProductsByStoreIdAndNames.
func (mr *MockProductRepositoryMockRecorder) GetProductsByStoreIdAndNames(ctx, storeID, names interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByStoreIdAndNames", reflect.TypeOf((*MockProductRepository)(nil).GetProductsByStoreIdAndNames), ctx, storeID, names)
}

// GetUnitOfMeasureByName mocks base method.
func (m *MockProductRepository) GetUnitOfMeasureByName(ctx context.Context, name string) (*entity.UnitOfMeasure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitOfMeasureByName", ctx, name)
	ret0, _ := ret[0].(*entity.UnitOfMeasure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitOfMeasureByName indicates an expected call of GetUnitOfMeasureByName.
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasureByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitOfMeasureByName", reflect.TypeOf((*MockProductRepository)(nil).GetUnitOfMeasureByName), ctx, name)
}

// GetUnitOfMeasureBySymbol mocks base method.
func (m *MockProductRepository) GetUnitOfMeasureBySymbol(ctx context.Context, symbol string) (*entity.UnitOfMeasure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitOfMeasureBySymbol", ctx, symbol)
	ret0, _ := ret[0].(*entity.UnitOfMeasure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitOfMeasureBySymbol indicates an expected call of GetUnitOfMeasureBySymbol.
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasureBySymbol(ctx, symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitOfMeasureBySymbol", reflect.TypeOf((*MockProductRepository)(nil).GetUnitOfMeasureBySymbol), ctx, symbol)
}

// GetUnitOfMeasures mocks base method.
func (m *MockProductRepository) GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) ([]*entity.UnitOfMeasure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitOfMeasures", ctx, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.UnitOfMeasure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitOfMeasures indicates an expected call of GetUnitOfMeasures.
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasures(ctx, isIncludeDeactivated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitOfMeasures", reflect.TypeOf((*MockProductRepository)(nil).GetUnitOfMeasures), ctx, isIncludeDeactivated)
}

// UpsertProductCategory mocks base method.
func (m *MockProductRepository) UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductCategory", ctx, prodCategory)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductCategory indicates an expected call of UpsertProductCategory.
func (mr *MockProductRepositoryMockRecorder) UpsertProductCategory(ctx, prodCategory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProductCategory", reflect.TypeOf((*MockProductRepository)(nil).UpsertProductCategory), ctx, prodCategory)
}

// UpsertProductType mocks base method.
func (m *MockProductRepository) UpsertProductType(ctx context.Context, prodType *entity.ProductType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductType", ctx, prodType)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductType indicates an expected call of UpsertProductType.
func (mr *MockProductRepositoryMockRecorder) UpsertProductType(ctx, prodType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProductType", reflect.TypeOf((*MockProductRepository)(nil).UpsertProductType), ctx, prodType)
}

// UpsertProducts mocks base method.
func (m *MockProductRepository) UpsertProducts(ctx context.Context, product []*entity.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProducts", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProducts indicates an expected call of UpsertProducts.
func (mr *MockProductRepositoryMockRecorder) UpsertProducts(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProducts", reflect.TypeOf((*MockProductRepository)(nil).UpsertProducts), ctx, product)
}

// UpsertUnitOfMeasure mocks base method.
func (m *MockProductRepository) UpsertUnitOfMeasure(ctx context.Context, uom *entity.UnitOfMeasure) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertUnitOfMeasure", ctx, uom)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertUnitOfMeasure indicates an expected call of UpsertUnitOfMeasure.
func (mr *MockProductRepositoryMockRecorder) UpsertUnitOfMeasure(ctx, uom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUnitOfMeasure", reflect.TypeOf((*MockProductRepository)(nil).UpsertUnitOfMeasure), ctx, uom)
}
