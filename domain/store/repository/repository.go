package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
)

// StoreServiceRepository defines the interface for store-related operations.
type StoreServiceRepository interface {
	// CreateStore creates a new store.
	CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error)

	// GetStore retrieves a store by its ID.
	GetStore(ctx context.Context, storeID string) (*entity.Store, error)

	// UpdateStore updates an existing store.
	UpdateStore(ctx context.Context, storeID string, update *entity.Store) (*entity.Store, error)

	// DeleteStore deletes a store by its ID.
	DeleteStore(ctx context.Context, storeID string) error

	// ListStores lists all stores.
	ListStores(ctx context.Context) ([]*entity.Store, error)
}

type Storage interface {
	UploadImage(ctx context.Context, image string) (string, error)
}
