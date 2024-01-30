package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	UpsertProducts(ctx context.Context, product []*entity.Product) error
	GetProductsByStoreId(ctx context.Context, storeID uuid.UUID) ([]*entity.Product, error)
	GetProductsByStoreIdAndNames(ctx context.Context, storeID uuid.UUID, names []string) ([]*entity.Product, error)
	UpsertUnitOfMeasure(ctx context.Context, uom *entity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType *entity.ProductType) error
}
