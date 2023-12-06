package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
)

type StoreInterface interface {
	GetAll(ctx context.Context) ([]*entity.Store, error)
}
