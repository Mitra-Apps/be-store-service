package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	UpsertProducts(ctx context.Context, product []*entity.Product) error
	GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *uuid.UUID, isIncludeDeactivated bool) ([]*entity.Product, error)
	GetProductsByStoreIdAndNames(ctx context.Context, storeID uuid.UUID, names []string) ([]*entity.Product, error)
	GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) ([]*entity.UnitOfMeasure, error)
	GetUnitOfMeasureByName(ctx context.Context, name string) (*entity.UnitOfMeasure, error)
	GetUnitOfMeasureBySymbol(ctx context.Context, symbol string) (*entity.UnitOfMeasure, error)
	GetProductCategoryByName(ctx context.Context, name string) (*entity.ProductCategory, error)
	GetProductCategories(ctx context.Context, isIncludeDeactivated bool) ([]*entity.ProductCategory, error)
	GetProductTypes(ctx context.Context, productCategoryID uuid.UUID, isIncludeDeactivated bool) ([]*entity.ProductType, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	UpsertUnitOfMeasure(ctx context.Context, uom *entity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType *entity.ProductType) error
}
