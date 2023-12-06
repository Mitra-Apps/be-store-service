package service

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
)

type Service struct {
	storeRepository repository.StoreInterface
}

func New(storeRepository repository.StoreInterface) *Service {
	return &Service{storeRepository: storeRepository}
}

type ServiceInterface interface {
	GetAll(ctx context.Context) ([]*entity.Store, error)
}
