package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProducts(ctx context.Context, product []entity.Product) error
	GetProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error)
}
