package service

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
)

func (s *Service) GetAll(ctx context.Context) ([]*entity.Store, error) {
	stores, err := s.storeRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return stores, nil
}
