package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	imageRepository "github.com/Mitra-Apps/be-store-service/domain/image/repository"
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
	GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) (products []*prodEntity.Product, err error)
	GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) (uom []*prodEntity.UnitOfMeasure, err error)
	GetProductCategories(ctx context.Context, isIncludeDeactivated bool) (cat []*prodEntity.ProductCategory, err error)
	GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) (types []*prodEntity.ProductType, err error)
	GetStoreByUserID(ctx context.Context, userID uuid.UUID) (store *entity.Store, err error)
}
type service struct {
	storeRepository   repository.StoreServiceRepository
	productRepository prodRepository.ProductRepository
	storage           repository.Storage
	imageRepository   imageRepository.ImageRepository
}

func New(
	storeRepository repository.StoreServiceRepository,
	prodRepository prodRepository.ProductRepository,
	storage repository.Storage,
	imageRepo imageRepository.ImageRepository,
) Service {
	return &service{
		storeRepository:   storeRepository,
		productRepository: prodRepository,
		storage:           storage,
		imageRepository:   imageRepo,
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
	if store == nil {
		return status.Errorf(codes.InvalidArgument, "Invalid store id")
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
		return err
	}

	var isAdmin bool
	for _, r := range roleNames {
		if r == "admin" {
			isAdmin = true
			break
		}
	}
	if existingStoreByStoreId.UserID != userID && !isAdmin {
		return status.Errorf(codes.PermissionDenied, "You don't have permission to create / update product for this store")
	}

	names := []string{}
	uomIds := []int64{}
	uomIdsMap := make(map[int64]bool)
	productTypeIds := []int64{}
	prodTypeIdsMap := make(map[int64]bool)
	for _, p := range products {
		p.StoreID = storeID
		names = append(names, p.Name)
		if p.UomID == 0 {
			return status.Errorf(codes.InvalidArgument, "Uom id is required")
		}
		if p.ProductTypeID == 0 {
			return status.Errorf(codes.InvalidArgument, "Product type id is required")
		}
		if !uomIdsMap[p.UomID] {
			uomIds = append(uomIds, p.UomID)
			uomIdsMap[p.UomID] = true
		}
		if !prodTypeIdsMap[p.ProductTypeID] {
			productTypeIds = append(productTypeIds, p.ProductTypeID)
			prodTypeIdsMap[p.ProductTypeID] = true
		}
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

	existingUoms, err := s.productRepository.GetUnitOfMeasuresByIds(ctx, uomIds)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when getting related uom")
	}
	if len(uomIds) > len(existingUoms) {
		return status.Errorf(codes.NotFound, "Unit of measure id is not found")
	}
	existingProdTypes, err := s.productRepository.GetProductTypesByIds(ctx, productTypeIds)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when getting related product type")
	}
	if len(productTypeIds) > len(existingProdTypes) {
		return status.Errorf(codes.NotFound, "Product type id is not found")
	}

	s.productRepository.InitiateTransaction(ctx)
	err = s.productRepository.UpsertProducts(ctx, products)
	if err != nil {
		return status.Errorf(codes.Internal, "Error when inserting products : "+err.Error())
	}
	var productImages []*prodEntity.ProductImage
	for _, p := range products {
		for _, i := range p.Images {
			if i != nil {
				imageId, err := s.imageRepository.UploadImage(ctx, i.ImageBase64Str, "product", userID.String())
				if err != nil {
					s.productRepository.TransactionRollback()
					return err
				}
				i.ImageId = *imageId
				productImages = append(productImages, i)
			}
		}
	}

	if err := s.productRepository.UpsertProductImages(ctx, productImages); err != nil {
		return err
	}

	if err := s.productRepository.TransactionCommit(); err != nil {
		return err
	}

	log.Println("Product is successfully inserted")

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
	return s.productRepository.UpsertProductCategory(ctx, prodCategory)
}

func (s *service) UpsertProductType(ctx context.Context, prodType *prodEntity.ProductType) error {
	if prodCat, err := s.productRepository.GetProductCategoryById(ctx, prodType.ProductCategoryID); err != nil {
		return status.Errorf(codes.AlreadyExists, "Error getting product category by id data")
	} else if prodCat == nil {
		return status.Errorf(codes.NotFound, "Related product category data is not found")
	}
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

func (s *service) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) (products []*prodEntity.Product, err error) {
	if _, err := s.storeRepository.GetStore(ctx, storeID.String()); err != nil {
		return nil, err
	}

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
	if prodCat, err := s.productRepository.GetProductCategoryById(ctx, productCategoryID); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Error getting product category by id data")
	} else if prodCat == nil {
		return nil, status.Errorf(codes.NotFound, "Product category id is not found")
	}

	if types, err = s.productRepository.GetProductTypes(ctx, productCategoryID, isIncludeDeactivated); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product types :"+err.Error())
	}
	return types, nil
}

func (s *service) GetProductById(ctx context.Context, id uuid.UUID) (p *prodEntity.Product, err error) {
	if p, err = s.productRepository.GetProductById(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product by id :"+err.Error())
	} else if p == nil && err == nil {
		return nil, status.Errorf(codes.NotFound, "Product id not found")
	}
	return p, nil
}

func (s *service) GetStoreByUserID(ctx context.Context, userID uuid.UUID) (store *entity.Store, err error) {
	store, err = s.storeRepository.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting store by user id :"+err.Error())
	}
	return store, nil
}
