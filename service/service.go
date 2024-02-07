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
	ListStores(ctx context.Context, page int32, limit int32) ([]*entity.Store, error)
	DeleteStores(ctx context.Context, storeIDs []string) error
	OpenCloseStore(ctx context.Context, userID uuid.UUID, roleNames []string, storeID string, isActive bool) error
	UpsertProducts(ctx context.Context, userID uuid.UUID, roleNames []string, storeID uuid.UUID, products []*prodEntity.Product) error
	UpsertUnitOfMeasure(ctx context.Context, uom *prodEntity.UnitOfMeasure) error
	UpsertProductCategory(ctx context.Context, prodCategory *prodEntity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType *prodEntity.ProductType) error
	GetProductById(ctx context.Context, id uuid.UUID) (*prodEntity.Product, error)
	GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *uuid.UUID, isIncludeDeactivated bool) (products []*prodEntity.Product, err error)
	GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) (uom []*prodEntity.UnitOfMeasure, err error)
	GetProductCategories(ctx context.Context, isIncludeDeactivated bool) (cat []*prodEntity.ProductCategory, err error)
	GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) (types []*prodEntity.ProductType, err error)
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

	for _, hour := range store.Hours {
		if hour.Is24Hr {
			hour.Open = "00:00"
			hour.Close = "23:59"
		}
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

	if claims.UserID.String() != update.UserID.String() || !claims.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to update this store")
	}

	update.UpdatedBy = claims.UserID
	update.ID, err = uuid.Parse(storeID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "store id should be uuid")
	}

	if exist, err := s.GetStore(ctx, storeID); err != nil {
		return nil, err
	} else if exist.UserID != claims.UserID {
		return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to update this store")
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

		if hour.Is24Hr {
			hour.Open = "00:00"
			hour.Close = "23:59"
		}
	}

	return s.storeRepository.UpdateStore(ctx, update)
}

func (s *service) ListStores(ctx context.Context, page int32, limit int32) ([]*entity.Store, error) {
	return s.storeRepository.ListStores(ctx, int(page), int(limit))
}

func (s *service) DeleteStores(ctx context.Context, storeIDs []string) error {
	return s.storeRepository.DeleteStores(ctx, storeIDs)
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
		p.StoreID = storeID
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
	existingUom, err := s.productRepository.GetUnitOfMeasureByName(ctx, uom.Name)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when getting uom by name : "+err.Error())
	}
	if existingUom != nil {
		return status.Errorf(codes.AlreadyExists, "Uom name is already exist in database")
	}
	existingUom, err = s.productRepository.GetUnitOfMeasureBySymbol(ctx, uom.Symbol)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when getting uom by symbol : "+err.Error())
	}
	if existingUom != nil {
		return status.Errorf(codes.AlreadyExists, "Uom symbol is already exist in database")
	}
	if err := s.productRepository.UpsertUnitOfMeasure(ctx, uom); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating unit of measure :"+err.Error())
	}
	return nil
}

func (s *service) UpsertProductCategory(ctx context.Context, prodCategory *prodEntity.ProductCategory) error {
	existingCat, err := s.productRepository.GetProductCategoryByName(ctx, prodCategory.Name)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when getting product category by name : "+err.Error())
	}
	if existingCat != nil {
		return status.Errorf(codes.AlreadyExists, "Category name is already exist in database")
	}

	if err := s.productRepository.UpsertProductCategory(ctx, prodCategory); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating product category :"+err.Error())
	}
	return nil
}

func (s *service) UpsertProductType(ctx context.Context, prodType *prodEntity.ProductType) error {
	existingProdType, err := s.productRepository.GetProductTypeByName(ctx, prodType.ProductCategoryID, prodType.Name)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when getting product type by name : "+err.Error())
	}
	if existingProdType != nil {
		return status.Errorf(codes.AlreadyExists, "Product type is already exist for this product category")
	}
	if err := s.productRepository.UpsertProductType(ctx, prodType); err != nil {
		return status.Errorf(codes.Internal, "Error when inserting / updating product type :"+err.Error())
	}
	return nil
}

func (s *service) GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) (uom []*prodEntity.UnitOfMeasure, err error) {
	if uom, err = s.productRepository.GetUnitOfMeasures(ctx, isIncludeDeactivated); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting unit of measures :"+err.Error())
	}
	return uom, nil
}

func (s *service) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *uuid.UUID, isIncludeDeactivated bool) (products []*prodEntity.Product, err error) {
	if products, err = s.productRepository.GetProductsByStoreId(ctx, storeID, productTypeId, isIncludeDeactivated); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product list :"+err.Error())
	}
	return products, nil
}
func (s *service) GetProductCategories(ctx context.Context, isIncludeDeactivated bool) (cat []*prodEntity.ProductCategory, err error) {
	if cat, err = s.productRepository.GetProductCategories(ctx, isIncludeDeactivated); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product categories :"+err.Error())
	}
	return cat, nil
}
func (s *service) GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) (types []*prodEntity.ProductType, err error) {
	if types, err = s.productRepository.GetProductTypes(ctx, productCategoryID, isIncludeDeactivated); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product types :"+err.Error())
	}
	return types, nil
}

func (s *service) GetProductById(ctx context.Context, id uuid.UUID) (p *prodEntity.Product, err error) {
	if p, err = s.productRepository.GetProductById(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product by id :"+err.Error())
	}
	return p, nil
}
