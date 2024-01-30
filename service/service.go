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
	UpdateStore(ctx context.Context, storeID string, update *entity.Store) (*entity.Store, error)
	GetStore(ctx context.Context, storeID string) (*entity.Store, error)
	ListStores(ctx context.Context) ([]*entity.Store, error)
	OpenCloseStore(ctx context.Context, userID uuid.UUID, roleNames []string, storeID string, isActive bool) error
	UpsertProducts(ctx context.Context, userID uuid.UUID, roleNames []string, storeID uuid.UUID, products []*prodEntity.Product) error
	UpsertUnitOfMeasure(ctx context.Context, uom *prodEntity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory *prodEntity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType *prodEntity.ProductType) error
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

func (s *service) UpdateStore(ctx context.Context, storeID string, update *entity.Store) (*entity.Store, error) {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Error when getting user id")
	}

	if strings.Compare(claims.UserID.String(), update.UserID.String()) != 0 {
		isAdmin := false
		for _, r := range claims.RoleNames {
			if r == "admin" {
				isAdmin = true
			}
		}

		if !isAdmin {
			return nil, status.Errorf(codes.PermissionDenied, "Only admin can update store")
		}
	}

	update.UpdatedBy = claims.UserID
	update.ID, err = uuid.Parse(storeID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "store id should be uuid")
	}

	for _, img := range update.Images {
		if img.ImageBase64 != "" {
			if img.ImageURL, err = s.storage.UploadImage(ctx, img.ImageBase64, storeID); err != nil {
				return nil, err
			}
		}

		img.UpdatedBy = claims.UserID
		img.ImageBase64 = ""
		img.StoreID = update.ID
		img.ID = uuid.Nil
	}

	for _, tag := range update.Tags {
		tag.UpdatedBy = claims.UserID
		tag.ID = uuid.Nil
	}

	for _, hour := range update.Hours {
		hour.UpdatedBy = claims.UserID
		hour.StoreID = update.ID
		hour.ID = uuid.Nil
	}

	return s.storeRepository.UpdateStore(ctx, update)
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

func (s *service) UpsertProducts(ctx context.Context, userID uuid.UUID, roleNames []string, storeID uuid.UUID, products []*prodEntity.Product) error {
	if len(products) == 0 {
		return status.Errorf(codes.InvalidArgument, "No product inserted")
	}
	existingStoreByStoreId, err := s.storeRepository.GetStore(ctx, storeID.String())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	var isAdmin bool
	for _, r := range roleNames {
		if r == "admin" {
			isAdmin = true
		}
	}
	if existingStoreByStoreId.UserID != userID && !isAdmin {
		return status.Errorf(codes.PermissionDenied, "You don't have permission to create / update product for this store")
	}
	names := []string{}
	for _, p := range products {
		names = append(names, p.Name)
	}
	existingProds, err := s.productRepository.GetProductsByStoreIdAndNames(ctx, storeID, names)
	if err != nil {
		return err
	}
	if len(existingProds) > 0 {
		existingProdNames := []string{}
		for _, p := range existingProds {
			existingProdNames = append(existingProdNames, p.Name)
		}
		return status.Errorf(codes.AlreadyExists,
			fmt.Sprintf("Product are already exist : %s", strings.Join(existingProdNames, ",")))
	}

	err = s.productRepository.UpsertProducts(ctx, products)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when inserting products : "+err.Error())
	}
	return nil
}

func (s *service) UpsertUnitOfMeasure(ctx context.Context, uom *prodEntity.UnitOfMeasure) error {
	if err := s.productRepository.UpsertUnitOfMeasure(ctx, uom); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating unit of measure :"+err.Error())
	}
	return nil
}

func (s *service) UpsertProductCategory(ctx context.Context, prodCategory *prodEntity.ProductCategory) error {
	if err := s.productRepository.UpsertProductCategory(ctx, prodCategory); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating product category :"+err.Error())
	}
	return nil
}

func (s *service) UpsertProductType(ctx context.Context, prodType *prodEntity.ProductType) error {
	if err := s.productRepository.UpsertProductType(ctx, prodType); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating product type :"+err.Error())
	}
	return nil
}
