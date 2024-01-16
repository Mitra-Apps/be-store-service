package service

import (
	"context"
	"strings"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error)
	GetStore(ctx context.Context, storeID string) (*entity.Store, error)
	ListStores(ctx context.Context) ([]*entity.Store, error)
	OpenCloseStore(ctx context.Context, storeID string, isActive bool) error
}
type service struct {
	storeRepository repository.StoreServiceRepository
	storage         repository.Storage
}

func New(
	storeRepository repository.StoreServiceRepository,
	storage repository.Storage,
) Service {
	return &service{
		storeRepository: storeRepository,
		storage:         storage,
	}
}

func (s *service) CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	for _, img := range store.Images {
		imageURL, err := s.storage.UploadImage(ctx, img.ImageBase64, store.UserID.String())
		if err != nil {
			return nil, err
		}
		img.ImageURL = imageURL
	}

	return s.storeRepository.CreateStore(ctx, store)
}

func (s *service) GetStore(ctx context.Context, storeID string) (*entity.Store, error) {
	return s.storeRepository.GetStore(ctx, storeID)
}

func (s *service) ListStores(ctx context.Context) ([]*entity.Store, error) {
	return s.storeRepository.ListStores(ctx)
}

func (s *service) OpenCloseStore(ctx context.Context, storeID string, isActive bool) error {
	if strings.Trim(storeID, " ") == "" {
		return status.Errorf(codes.InvalidArgument, "store id is required")
	}
	storeIDUuid, err := uuid.Parse(storeID)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "store id should be uuid")
	}

	err = s.storeRepository.OpenCloseStore(ctx, storeIDUuid, isActive)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when opening / closing store")
	}
	return nil
}
