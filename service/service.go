package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	imageRepository "github.com/Mitra-Apps/be-store-service/domain/image/repository"
	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	prodRepository "github.com/Mitra-Apps/be-store-service/domain/product/repository"
	errPb "github.com/Mitra-Apps/be-store-service/domain/proto"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	"github.com/Mitra-Apps/be-store-service/handler/grpc/middleware"
	"github.com/Mitra-Apps/be-store-service/lib"
	utilityPb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	util "github.com/Mitra-Apps/be-utility-service/service"
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
	UpsertProducts(ctx context.Context, userID uuid.UUID, roleNames []string, storeID uuid.UUID, isUpdate bool, products ...*prodEntity.Product) error
	UpsertProductCategory(ctx context.Context, prodCategory *prodEntity.ProductCategory) error
	UpsertProductType(ctx context.Context, prodType *prodEntity.ProductType) error
	GetProductById(ctx context.Context, id uuid.UUID) (*prodEntity.Product, error)
	GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) (products []*prodEntity.Product, err error)
	DeleteProductById(ctx context.Context, userId uuid.UUID, id uuid.UUID) error
	GetProductCategories(ctx context.Context, isIncludeDeactivated bool) (cat []*prodEntity.ProductCategory, uom []string, err error)
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

		if !hour.IsOpen {
			hour.Open = "00:00"
			hour.Close = "00:00"
		} else if hour.Is24Hr {
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

	if claims.UserID.String() != update.UserID.String() && !claims.IsAdmin {
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

		if !hour.IsOpen {
			hour.Open = "00:00"
			hour.Close = "00:00"
		} else if hour.Is24Hr {
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

func (s *service) UpsertProducts(ctx context.Context, userID uuid.UUID, roleNames []string, storeID uuid.UUID, isUpdate bool, products ...*prodEntity.Product) error {
	if len(products) == 0 {
		return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_NO_PRODUCT_INSERTED.String(), "tidak ada produk yang disimpan")
	}

	if storeID == uuid.Nil {
		return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_STORE_ID_IS_REQUIRED.String(), "Store id diperlukan")
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
		return util.NewError(codes.PermissionDenied, errPb.StoreErrorCode_DONT_HAVE_PERMISSION_TO_CREATE_OR_UPDATE_STORE.String(), "Anda tidak memiliki izin untuk membuat/memperbarui produk untuk toko ini")
	}

	names := []string{}
	productTypeIds := []int64{}
	prodTypeIdsMap := make(map[int64]bool)
	for _, p := range products {
		if isUpdate && p.ID == uuid.Nil {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_IS_REQUIRED.String(), "Product id diperlukan")
		} else if !isUpdate && p.ID != uuid.Nil {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_ID_SHOULD_BE_EMPTY.String(), "Product id harus kosong")
		}

		p.StoreID = storeID
		names = append(names, p.Name)
		if p.Uom == "" {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_UOM_IS_REQUIRED.String(), "Satuan unit diperlukan")
		}
		if p.ProductTypeID == 0 {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_TYPE_IS_REQUIRED.String(), "Product type id diperlukan")
		}
		if p.Stock < 0 {
			return util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_STOCK_SHOULD_BE_POSITIVE.String(), "Stock harus positif")
		}
		if !prodTypeIdsMap[p.ProductTypeID] {
			productTypeIds = append(productTypeIds, p.ProductTypeID)
			prodTypeIdsMap[p.ProductTypeID] = true
		}
	}

	if !isUpdate {
		existingProds, err := s.productRepository.GetProductsByStoreIdAndNames(ctx, storeID, names)
		if err != nil {
			return err
		}
		if len(existingProds) > 0 {
			existingProdNames := []string{}
			for _, p := range existingProds {
				existingProdNames = append(existingProdNames, p.Name)
			}
			return util.NewError(codes.AlreadyExists, errPb.StoreErrorCode_PRODUCTS_ARE_ALREADY_REGISTERED.String(),
				fmt.Sprintf("Produk sudah terdaftar : %s", strings.Join(existingProdNames, ",")))
		}
	}

	existingProdTypes, err := s.productRepository.GetProductTypesByIds(ctx, productTypeIds)
	if err != nil {
		return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_GETTING_RELATED_PRODUCT_TYPE.String(), "Error saat mencari data product type")
	}

	if len(productTypeIds) > len(existingProdTypes) {
		return util.NewError(codes.NotFound, errPb.StoreErrorCode_PRODUCT_TYPE_ID_IS_NOT_FOUND.String(), "Product type id tidak ditemukan")
	}

	s.productRepository.InitiateTransaction(ctx)
	err = s.productRepository.UpsertProducts(ctx, products)
	if err != nil {
		return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_INSERTING_OR_UPDATING_PRODUCT.String(), "Error saat membuat / memperbarui data produk : "+err.Error())
	}

	var prodIds []uuid.UUID
	for _, p := range products {
		prodIds = append(prodIds, p.ID)
	}

	var addProductImages []*prodEntity.ProductImage
	if isUpdate {
		// add product images if the product image id is nil
		// remove product images if if the product image id exist in db but not exist in the request.
		// if the product images id exist in the db and in the request, do nothing.
		var removeProductImages []*prodEntity.ProductImage

		prodImages, existingProdImagesByProdIdMap, err := s.productRepository.GetProductImagesByProductIds(ctx, prodIds)
		if err != nil {
			return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_GETTING_PRODUCT_IMAGE.String(), "Error saat melakukan pencarian gambar produk : "+err.Error())
		}

		existingProdImagesMap := make(map[uuid.UUID]*prodEntity.ProductImage)
		for _, i := range prodImages {
			existingProdImagesMap[i.ID] = i
		}

		for _, p := range products {
			fromReqProdImagesMap := make(map[uuid.UUID]*prodEntity.ProductImage)
			for _, i := range p.Images {
				if i.BaseModel.ID == uuid.Nil {
					i.ProductId = p.ID
					addProductImages = append(addProductImages, i)
				} else {
					fromReqProdImagesMap[i.ID] = i
				}
			}

			// remove product images if if the product image id exist in db but not exist in the request.
			for _, i := range existingProdImagesByProdIdMap[p.ID] {
				if fromReqProdImagesMap[i.ID] == nil {
					removeProductImages = append(removeProductImages, i)
				}
			}

			// upload new images
			for _, i := range addProductImages {
				imageId, err := s.UploadImageToStorage(ctx, i.ImageBase64Str, userID)
				if err != nil {
					s.productRepository.TransactionRollback()
					return err
				}
				i.ImageId = *imageId
			}

			// remove images from storage and remove product_image data from database
			if len(removeProductImages) > 0 {
				removeProdImageIds := []string{}
				for _, i := range removeProductImages {
					log.Println("Remove : ", i.ID.String())
					removeProdImageIds = append(removeProdImageIds, i.ID.String())
				}
				log.Printf("Remove product images : %v \n", removeProdImageIds)
				if err := s.imageRepository.RemoveImage(ctx, removeProdImageIds, "product", userID.String()); err != nil {
					s.productRepository.TransactionRollback()
					return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_REMOVING_IMAGE_FROM_STORAGE.String(), "Error saat menghapus gambar dari penyimpanan : "+err.Error())
				}
				if err := s.productRepository.DeleteProductImages(ctx, removeProductImages); err != nil {
					s.productRepository.TransactionRollback()
					return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_DELETING_PRODUCT_IMAGE.String(), "Error saat menghapus gambar produk : "+err.Error())
				}
			}
		}
	} else {
		for _, p := range products {
			for _, i := range p.Images {
				if i != nil {
					imageId, err := s.UploadImageToStorage(ctx, i.ImageBase64Str, userID)
					if err != nil {
						return err
					}
					i.ProductId = p.ID
					i.ImageId = *imageId
					addProductImages = append(addProductImages, i)
				}
			}
		}
	}

	if len(addProductImages) > 0 {
		log.Printf("Add product images count : %d \n", len(addProductImages))
		if err := s.productRepository.UpsertProductImages(ctx, addProductImages); err != nil {
			s.productRepository.TransactionRollback()
			return err
		}
	}

	if err := s.productRepository.TransactionCommit(); err != nil {
		return err
	}

	log.Println("Product is successfully inserted")

	return nil
}

func (s *service) UploadImageToStorage(ctx context.Context, imageBase64Str string, userID uuid.UUID) (*uuid.UUID, error) {
	if strings.Trim(imageBase64Str, " ") == "" {
		return nil, util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_IMAGE_SHOULD_BE_IN_BASE_64_FORMAT.String(), "Gambar produk harus dalam format base 64")
	}

	imageId, err := s.imageRepository.UploadImage(ctx, imageBase64Str, "product", userID.String())
	if err != nil {
		s.productRepository.TransactionRollback()
		return nil, err
	}

	return imageId, nil
}

func (s *service) UpsertProductCategory(ctx context.Context, prodCategory *prodEntity.ProductCategory) error {
	existingProdCategory, err := s.productRepository.GetProductCategoryByName(ctx, strings.ToLower(prodCategory.Name))
	if err != nil {
		return err
	}

	if existingProdCategory != nil && existingProdCategory.ID != prodCategory.ID {
		return status.Errorf(codes.AlreadyExists, "Name is already used by another product category")
	}

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
		return status.Errorf(codes.Internal, "Error when inserting / updating product Type :"+err.Error())
	}

	return nil
}

func (s *service) DeleteProductById(ctx context.Context, userId uuid.UUID, id uuid.UUID) error {
	prodIds := []uuid.UUID{id}
	prodImages, _, err := s.productRepository.GetProductImagesByProductIds(ctx, prodIds)
	if err != nil {
		return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_GETTING_PRODUCT_IMAGE.String(), "Error saat melakukan pencarian gambar produk : "+err.Error())
	}

	s.productRepository.InitiateTransaction(ctx)
	removeProdImageIds := []string{}
	for _, i := range prodImages {
		log.Println("Remove : ", i.ID.String())
		removeProdImageIds = append(removeProdImageIds, i.ID.String())
	}

	log.Printf("Remove product images : %v \n", removeProdImageIds)
	if len(removeProdImageIds) > 0 {
		if err := s.imageRepository.RemoveImage(ctx, removeProdImageIds, "product", userId.String()); err != nil {
			s.productRepository.TransactionRollback()
			return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_REMOVING_IMAGE_FROM_STORAGE.String(), "Error saat menghapus gambar dari penyimpanan : "+err.Error())
		}

		if err := s.productRepository.DeleteProductImages(ctx, prodImages); err != nil {
			s.productRepository.TransactionRollback()
			return util.NewError(codes.Internal, errPb.StoreErrorCode_ERROR_WHEN_DELETING_PRODUCT_IMAGE.String(), "Error saat menghapus gambar produk : "+err.Error())
		}
	}

	err = s.productRepository.DeleteProductById(ctx, id)
	if err != nil {
		s.productRepository.TransactionRollback()
		return util.NewError(codes.Internal, errPb.StoreErrorCode_PRODUCT_IS_REQUIRED.String(), "Error saat menghapus produk : "+err.Error())
	}

	s.productRepository.TransactionCommit()

	return err
}

func (s *service) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) (products []*prodEntity.Product, err error) {
	if _, err := s.storeRepository.GetStore(ctx, storeID.String()); err != nil {
		return nil, err
	}

	if products, err = s.productRepository.GetProductsByStoreId(ctx, storeID, productTypeId, isIncludeDeactivated); err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting product list :"+err.Error())
	}

	if err := s.GetProductImagesInformation(ctx, nil, products); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *service) GetProductCategories(ctx context.Context, isIncludeDeactivated bool) (cat []*prodEntity.ProductCategory, uom []string, err error) {
	if cat, err = s.productRepository.GetProductCategories(ctx, isIncludeDeactivated); err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Error when getting product categories :"+err.Error())
	}

	// read file uom
	fileName := "uom.json"
	var uoms []string
	err = lib.ReadToFile(fileName, lib.JsonFormat, &uoms)
	if err != nil {
		err = nil
		// return nil, nil, status.Errorf(codes.Internal, err.Error())
	}

	return cat, uoms, nil
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

	if err := s.GetProductImagesInformation(ctx, p, nil); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetProductImagesInformation(ctx context.Context, product *prodEntity.Product, products []*prodEntity.Product) error {
	prodImages := []*prodEntity.ProductImage{}
	if product == nil {
		for _, p := range products {
			prodImages = append(prodImages, p.Images...)
		}
	} else {
		prodImages = append(prodImages, product.Images...)
	}

	if len(prodImages) > 0 {
		ids := []string{}
		for _, img := range prodImages {
			ids = append(ids, img.ImageId.String())
		}

		images, err := s.imageRepository.GetImagesByIds(ctx, ids)
		if err != nil {
			return err
		}

		imgMap := make(map[string]*utilityPb.Image)
		for _, i := range images {
			imgMap[i.Id] = i
		}

		if product == nil {
			for _, p := range products {
				for _, img := range p.Images {
					img.ImageURL = imgMap[img.ImageId.String()].Path
				}
			}
		} else {
			for _, img := range product.Images {
				img.ImageURL = imgMap[img.ImageId.String()].Path
			}
		}
	}

	return nil
}

func (s *service) GetStoreByUserID(ctx context.Context, userID uuid.UUID) (store *entity.Store, err error) {
	store, err = s.storeRepository.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error when getting store by user id :"+err.Error())
	}

	return store, nil
}
