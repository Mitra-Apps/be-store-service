package repository

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/google/uuid"
)

// StoreServiceRepository defines the interface for store-related operations.
type StoreServiceRepository interface {
	// CreateStore creates a new store.
	CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error)

	// GetStore retrieves a store by its ID.
	GetStore(ctx context.Context, storeID string) (*entity.Store, error)

	// UpdateStore updates an existing store.
	UpdateStore(ctx context.Context, update *entity.Store) (*entity.Store, error)

	// DeleteStore deletes a store by its ID.
	DeleteStore(ctx context.Context, storeID string) error

	// ListStores lists all stores.
	ListStores(ctx context.Context, page, pageSize int) ([]*entity.Store, error)

	// Activate / Deactivate store
	OpenCloseStore(ctx context.Context, storeID uuid.UUID, isActive bool) error

	// GetStoreByUserID retrieves a store by its user ID.
	GetStoreByUserID(ctx context.Context, userID uuid.UUID) (*entity.Store, error)
}

type Storage interface {
	UploadImage(ctx context.Context, image, userID string) (string, error)
}
