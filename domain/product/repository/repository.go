package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	UpsertProducts(ctx context.Context, product []entity.Product) error
	GetProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error)
	UpsertUnitOfMeasure(ctx context.Context, uom entity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory entity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType entity.ProductType) error
}
