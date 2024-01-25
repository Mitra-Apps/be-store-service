package service

import (
	"context"
	"fmt"
	"strings"

	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	prodRepository "github.com/Mitra-Apps/be-store-service/domain/product/repository"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	"github.com/Mitra-Apps/be-store-service/handler/grpc/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error)
	GetStore(ctx context.Context, storeID string) (*entity.Store, error)
	ListStores(ctx context.Context) ([]*entity.Store, error)
	OpenCloseStore(ctx context.Context, userID uuid.UUID, roleNames []string, storeID string, isActive bool) error
	UpsertProducts(ctx context.Context, products []prodEntity.Product) error
	UpsertUnitOfMeasure(ctx context.Context, uom prodEntity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory prodEntity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType prodEntity.ProductType) error
}
type service struct {
	storeRepository   repository.StoreServiceRepository
	productRepository prodRepository.ProductRepository
	storage           repository.Storage
}

func New(
	storeRepository repository.StoreServiceRepository,
	prodRepository prodRepository.ProductRepository,
	storage repository.Storage,
) Service {
	return &service{
		storeRepository:   storeRepository,
		productRepository: prodRepository,
		storage:           storage,
	}
}

func (s *service) CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error when getting user id")
	}

	store.UserID = claims.UserID
	store.CreatedBy = claims.UserID

	exist, err := s.storeRepository.GetStoreByUserID(ctx, store.UserID)
	if err != nil {
		return nil, err
	}

	if exist != nil {
		return nil, status.Errorf(codes.AlreadyExists, "User already has a store")
	}

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

func (s *service) OpenCloseStore(ctx context.Context, userID uuid.UUID, roleNames []string, storeID string, isActive bool) error {
	if strings.Trim(storeID, " ") == "" {
		return status.Errorf(codes.InvalidArgument, "store id is required")
	}
	storeIDUuid, err := uuid.Parse(storeID)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "store id should be uuid")
	}

	store, err := s.storeRepository.GetStore(ctx, storeID)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	var isAdmin bool
	for _, r := range roleNames {
		if r == "admin" {
			isAdmin = true
		}
	}

	if userID != store.UserID && !isAdmin {
		return status.Errorf(codes.PermissionDenied, "You do not have permission to open / close this store")
	}

	err = s.storeRepository.OpenCloseStore(ctx, storeIDUuid, isActive)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when opening / closing store")
	}
	return nil
}

func (s *service) UpsertProducts(ctx context.Context, products []prodEntity.Product) error {
	if len(products) == 0 {
		return status.Errorf(codes.InvalidArgument, "No product inserted")
	}
	ids := []uuid.UUID{}
	for _, p := range products {
		ids = append(ids, p.ID)
	}
	existingProds, err := s.productRepository.GetProductsByIds(ctx, ids)
	if err != nil {
		return err
	}
	if len(existingProds) > 0 {
		existingProdIds := []string{}
		for _, p := range existingProds {
			existingProdIds = append(existingProdIds, p.ID.String())
		}
		return status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Product are already exist : %s", strings.Join(existingProdIds, ",")))
	}

	err = s.productRepository.UpsertProducts(ctx, products)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when inserting products")
	}
	return nil
}

func (s *service) UpsertUnitOfMeasure(ctx context.Context, uom prodEntity.UnitOfMeasure) error {
	if err := s.productRepository.UpsertUnitOfMeasure(ctx, uom); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating unit of measure")
	}
	return nil
}

func (s *service) UpsertProductCategory(ctx context.Context, prodCategory prodEntity.ProductCategory) error {
	if err := s.productRepository.UpsertProductCategory(ctx, prodCategory); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating product category")
	}
	return nil
}

func (s *service) UpsertProductType(ctx context.Context, prodType prodEntity.ProductType) error {
	if err := s.productRepository.UpsertProductType(ctx, prodType); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating product type")
	}
	return nil
}
