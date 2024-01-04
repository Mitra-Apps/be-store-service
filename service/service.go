package service

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
)

type Service interface {
	CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error)
	GetStore(ctx context.Context, storeID string) (*entity.Store, error)
	ListStores(ctx context.Context) ([]*entity.Store, error)
}
type service struct {
	storeRepository repository.StoreServiceRepository
}

func New(
	storeRepository repository.StoreServiceRepository,
) Service {
	return &service{
		storeRepository: storeRepository,
	}
}

func (s *service) CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	return s.storeRepository.CreateStore(ctx, store)
}

func (s *service) GetStore(ctx context.Context, storeID string) (*entity.Store, error) {
	return s.storeRepository.GetStore(ctx, storeID)
}

func (s *service) ListStores(ctx context.Context) ([]*entity.Store, error) {
	return s.storeRepository.ListStores(ctx)
}
