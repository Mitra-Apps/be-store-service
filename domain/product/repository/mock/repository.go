// Code generated by MockGen. DO NOT EDIT.
// Source: domain/product/repository/repository.go
//
// Generated by this command:
//
//	mockgen -source=domain/product/repository/repository.go -destination=domain/product/repository/mock/repository.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
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
func (mr *MockProductRepositoryMockRecorder) GetProductById(ctx, id any) *gomock.Call {
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
func (mr *MockProductRepositoryMockRecorder) GetProductCategories(ctx, isIncludeDeactivated any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductCategories", reflect.TypeOf((*MockProductRepository)(nil).GetProductCategories), ctx, isIncludeDeactivated)
}

// GetProductCategoryById mocks base method.
func (m *MockProductRepository) GetProductCategoryById(ctx context.Context, id int64) (*entity.ProductCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductCategoryById", ctx, id)
	ret0, _ := ret[0].(*entity.ProductCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductCategoryById indicates an expected call of GetProductCategoryById.
func (mr *MockProductRepositoryMockRecorder) GetProductCategoryById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductCategoryById", reflect.TypeOf((*MockProductRepository)(nil).GetProductCategoryById), ctx, id)
}

// GetProductCategoryByName mocks base method.
func (m *MockProductRepository) GetProductCategoryByName(ctx context.Context, name string) (*entity.ProductCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductCategoryByName", ctx, name)
	ret0, _ := ret[0].(*entity.ProductCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductCategoryByName indicates an expected call of GetProductCategoryByName.
func (mr *MockProductRepositoryMockRecorder) GetProductCategoryByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductCategoryByName", reflect.TypeOf((*MockProductRepository)(nil).GetProductCategoryByName), ctx, name)
}

// GetProductTypeByName mocks base method.
func (m *MockProductRepository) GetProductTypeByName(ctx context.Context, productCategoryID int64, name string) (*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTypeByName", ctx, productCategoryID, name)
	ret0, _ := ret[0].(*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTypeByName indicates an expected call of GetProductTypeByName.
func (mr *MockProductRepositoryMockRecorder) GetProductTypeByName(ctx, productCategoryID, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTypeByName", reflect.TypeOf((*MockProductRepository)(nil).GetProductTypeByName), ctx, productCategoryID, name)
}

// GetProductTypes mocks base method.
func (m *MockProductRepository) GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) ([]*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTypes", ctx, productCategoryID, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTypes indicates an expected call of GetProductTypes.
func (mr *MockProductRepositoryMockRecorder) GetProductTypes(ctx, productCategoryID, isIncludeDeactivated any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTypes", reflect.TypeOf((*MockProductRepository)(nil).GetProductTypes), ctx, productCategoryID, isIncludeDeactivated)
}

// GetProductTypesByIds mocks base method.
func (m *MockProductRepository) GetProductTypesByIds(ctx context.Context, typeIds []int64) ([]*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTypesByIds", ctx, typeIds)
	ret0, _ := ret[0].([]*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTypesByIds indicates an expected call of GetProductTypesByIds.
func (mr *MockProductRepositoryMockRecorder) GetProductTypesByIds(ctx, typeIds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTypesByIds", reflect.TypeOf((*MockProductRepository)(nil).GetProductTypesByIds), ctx, typeIds)
}

// GetProductsByStoreId mocks base method.
func (m *MockProductRepository) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) ([]*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByStoreId", ctx, storeID, productTypeId, isIncludeDeactivated)
	ret0, _ := ret[0].([]*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsByStoreId indicates an expected call of GetProductsByStoreId.
func (mr *MockProductRepositoryMockRecorder) GetProductsByStoreId(ctx, storeID, productTypeId, isIncludeDeactivated any) *gomock.Call {
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
func (mr *MockProductRepositoryMockRecorder) GetProductsByStoreIdAndNames(ctx, storeID, names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByStoreIdAndNames", reflect.TypeOf((*MockProductRepository)(nil).GetProductsByStoreIdAndNames), ctx, storeID, names)
}

// GetUnitOfMeasureById mocks base method.
func (m *MockProductRepository) GetUnitOfMeasureById(ctx context.Context, uomId int64) (*entity.UnitOfMeasure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitOfMeasureById", ctx, uomId)
	ret0, _ := ret[0].(*entity.UnitOfMeasure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitOfMeasureById indicates an expected call of GetUnitOfMeasureById.
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasureById(ctx, uomId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitOfMeasureById", reflect.TypeOf((*MockProductRepository)(nil).GetUnitOfMeasureById), ctx, uomId)
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
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasureByName(ctx, name any) *gomock.Call {
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
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasureBySymbol(ctx, symbol any) *gomock.Call {
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
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasures(ctx, isIncludeDeactivated any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitOfMeasures", reflect.TypeOf((*MockProductRepository)(nil).GetUnitOfMeasures), ctx, isIncludeDeactivated)
}

// GetUnitOfMeasuresByIds mocks base method.
func (m *MockProductRepository) GetUnitOfMeasuresByIds(ctx context.Context, uomIds []int64) ([]*entity.UnitOfMeasure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitOfMeasuresByIds", ctx, uomIds)
	ret0, _ := ret[0].([]*entity.UnitOfMeasure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitOfMeasuresByIds indicates an expected call of GetUnitOfMeasuresByIds.
func (mr *MockProductRepositoryMockRecorder) GetUnitOfMeasuresByIds(ctx, uomIds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitOfMeasuresByIds", reflect.TypeOf((*MockProductRepository)(nil).GetUnitOfMeasuresByIds), ctx, uomIds)
}

// InitiateTransaction mocks base method.
func (m *MockProductRepository) InitiateTransaction(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitiateTransaction", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// InitiateTransaction indicates an expected call of InitiateTransaction.
func (mr *MockProductRepositoryMockRecorder) InitiateTransaction(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitiateTransaction", reflect.TypeOf((*MockProductRepository)(nil).InitiateTransaction), ctx)
}

// TransactionCommit mocks base method.
func (m *MockProductRepository) TransactionCommit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionCommit")
	ret0, _ := ret[0].(error)
	return ret0
}

// TransactionCommit indicates an expected call of TransactionCommit.
func (mr *MockProductRepositoryMockRecorder) TransactionCommit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionCommit", reflect.TypeOf((*MockProductRepository)(nil).TransactionCommit))
}

// TransactionRollback mocks base method.
func (m *MockProductRepository) TransactionRollback() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TransactionRollback")
}

// TransactionRollback indicates an expected call of TransactionRollback.
func (mr *MockProductRepositoryMockRecorder) TransactionRollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionRollback", reflect.TypeOf((*MockProductRepository)(nil).TransactionRollback))
}

// UpsertProductCategory mocks base method.
func (m *MockProductRepository) UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductCategory", ctx, prodCategory)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductCategory indicates an expected call of UpsertProductCategory.
func (mr *MockProductRepositoryMockRecorder) UpsertProductCategory(ctx, prodCategory any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProductCategory", reflect.TypeOf((*MockProductRepository)(nil).UpsertProductCategory), ctx, prodCategory)
}

// UpsertProductImages mocks base method.
func (m *MockProductRepository) UpsertProductImages(ctx context.Context, productImages []*entity.ProductImage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductImages", ctx, productImages)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductImages indicates an expected call of UpsertProductImages.
func (mr *MockProductRepositoryMockRecorder) UpsertProductImages(ctx, productImages any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProductImages", reflect.TypeOf((*MockProductRepository)(nil).UpsertProductImages), ctx, productImages)
}

// UpsertProductType mocks base method.
func (m *MockProductRepository) UpsertProductType(ctx context.Context, prodType *entity.ProductType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProductType", ctx, prodType)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProductType indicates an expected call of UpsertProductType.
func (mr *MockProductRepositoryMockRecorder) UpsertProductType(ctx, prodType any) *gomock.Call {
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
func (mr *MockProductRepositoryMockRecorder) UpsertProducts(ctx, product any) *gomock.Call {
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
func (mr *MockProductRepositoryMockRecorder) UpsertUnitOfMeasure(ctx, uom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUnitOfMeasure", reflect.TypeOf((*MockProductRepository)(nil).UpsertUnitOfMeasure), ctx, uom)
}
