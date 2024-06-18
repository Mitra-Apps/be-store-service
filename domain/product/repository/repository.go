package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/Mitra-Apps/be-store-service/types"
	"github.com/google/uuid"
)

type ProductRepository interface {
	UpsertProducts(ctx context.Context, product []*entity.Product) error
	GetProductsByStoreId(ctx context.Context, params types.GetProductsByStoreIdRepoParams) ([]*entity.Product, base_model.Pagination, error)
	GetProductsByStoreIdAndNames(ctx context.Context, storeID uuid.UUID, names []string) ([]*entity.Product, error)
	GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) ([]*entity.UnitOfMeasure, error)
	GetUnitOfMeasureByName(ctx context.Context, name string) (*entity.UnitOfMeasure, error)
	GetUnitOfMeasureBySymbol(ctx context.Context, symbol string) (*entity.UnitOfMeasure, error)
	GetUnitOfMeasuresByIds(ctx context.Context, uomIds []int64) ([]*entity.UnitOfMeasure, error)
	GetUnitOfMeasureById(ctx context.Context, uomId int64) (*entity.UnitOfMeasure, error)
	GetProductCategoryByName(ctx context.Context, name string) (*entity.ProductCategory, error)
	GetProductCategoryById(ctx context.Context, id int64) (*entity.ProductCategory, error)
	GetProductCategories(ctx context.Context, isIncludeDeactivated bool) ([]*entity.ProductCategory, error)
	GetProductCategoriesByStoreId(ctx context.Context, params types.GetProductCategoriesByStoreIdParams) ([]*entity.ProductCategory, error)
	GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) ([]*entity.ProductType, error)
	GetProductTypeByName(ctx context.Context, productCategoryID int64, name string) (*entity.ProductType, error)
	GetProductTypesByIds(ctx context.Context, typeIds []int64) ([]*entity.ProductType, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	UpsertUnitOfMeasure(ctx context.Context, uom *entity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType *entity.ProductType) error
	UpsertProductImages(ctx context.Context, productImages []*entity.ProductImage) error
	InitiateTransaction(ctx context.Context) bool
	TransactionCommit() error
	TransactionRollback()
	GetProductImagesByProductIds(ctx context.Context, productIds []uuid.UUID) ([]*entity.ProductImage, map[uuid.UUID][]*entity.ProductImage, error)
	DeleteProductImages(ctx context.Context, productImages []*entity.ProductImage) error
	DeleteProductById(ctx context.Context, id uuid.UUID) error
	GetProductTypeById(ctx context.Context, id int64) (*entity.ProductType, error)
	GetCategoriesByStoreIdRepo(ctx context.Context, input types.GetCategoriesByStoreIdRepoInput) (output types.GetCategoriesByStoreIdRepoOutput, err error)
}
