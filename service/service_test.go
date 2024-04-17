package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	imageRepoMock "github.com/Mitra-Apps/be-store-service/domain/image/repository/mock"
	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	prodRepoMock "github.com/Mitra-Apps/be-store-service/domain/product/repository/mock"
	errPb "github.com/Mitra-Apps/be-store-service/domain/proto"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	storeRepoMock "github.com/Mitra-Apps/be-store-service/domain/store/repository/mock"
	util "github.com/Mitra-Apps/be-utility-service/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	userID               = "8b15140c-f6d0-4f2f-8302-57383a51adaf"
	otherUserID          = "2f27d467-9f83-4170-96ab-36e0994f37ca"
	storeID              = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
	otherStoreID         = "52d11042-8c45-453e-86af-fe1e4d7facf6"
	otherStoreID2        = "52d11042-8c45-453e-86af-fe1e4d7facf7"
	productID            = "7d56be32-70a2-4f49-b66b-63e6f8e719d7"
	otherProductID       = "7d56be32-70a2-4f49-b66b-63e6f8e719d8"
	otherProductID2      = "7d56be32-70a2-4f49-b66b-63e6f8e719d9"
	productImageID       = "7d56be32-70a2-4f49-b66b-63e6f8e719e7"
	otherproductImageID2 = "7d56be32-70a2-4f49-b66b-63e6f8e719e9"
)

func TestOpenCloseStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)

	ctx := context.Background()

	service := New(mockStoreRepo, nil, nil, nil)

	userIdUuid, _ := uuid.Parse(userID)
	otherUserIdUuid, _ := uuid.Parse(otherUserID)
	storeIdUuid, _ := uuid.Parse(storeID)
	roleNames := []string{}
	roleNames = append(roleNames, "merchant")
	adminRole := []string{}
	adminRole = append(adminRole, "admin")
	store := &entity.Store{
		UserID: userIdUuid,
	}

	t.Run("Should return success", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		mockStoreRepo.EXPECT().
			OpenCloseStore(gomock.Any(), storeIdUuid, false).
			Times(1).
			Return(nil)

		err := service.OpenCloseStore(ctx, userIdUuid, roleNames, storeID, false)

		assert.Nil(t, err)
	})

	t.Run("Should return error if failed to get store", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(nil, errors.New("error"))

		err := service.OpenCloseStore(ctx, userIdUuid, roleNames, storeID, false)

		assert.Error(t, err)
	})

	t.Run("Should return error if store is nill", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(nil, nil)

		err := service.OpenCloseStore(ctx, userIdUuid, roleNames, storeID, false)

		assert.Error(t, err)
	})

	t.Run("Should return error if admin user id is different", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		mockStoreRepo.EXPECT().
			OpenCloseStore(gomock.Any(), storeIdUuid, false).
			Times(1).
			Return(nil)

		err := service.OpenCloseStore(ctx, otherUserIdUuid, adminRole, storeID, false)

		assert.Nil(t, err)
	})

	t.Run("Should return error if user id is not admin", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.OpenCloseStore(ctx, otherUserIdUuid, roleNames, storeID, false)

		assert.Error(t, err)
	})

	t.Run("Should return error if store id invalid", func(t *testing.T) {
		err := service.OpenCloseStore(ctx, userIdUuid, roleNames, "aaa", false)

		assert.Error(t, err)
	})

	t.Run("Should return error if store id not provided", func(t *testing.T) {
		err := service.OpenCloseStore(ctx, userIdUuid, roleNames, "", false)

		assert.Error(t, err)
	})

	t.Run("Should return error to open and close the store", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		mockStoreRepo.EXPECT().
			OpenCloseStore(gomock.Any(), storeIdUuid, false).
			Times(1).
			Return(errors.New("error"))

		err := service.OpenCloseStore(ctx, userIdUuid, roleNames, storeID, false)

		assert.Error(t, err)
	})
}

func TestCreateStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	storeRepository := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	storage := storeRepoMock.NewMockStorage(ctrl)
	service := New(storeRepository, nil, storage, nil)

	ctx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	sessionUserID := uuid.New()

	md.Set("x-user-id", sessionUserID.String())
	ctx = metadata.NewIncomingContext(ctx, md)

	imgBase64 := "SampleImageBase64"
	imgURL := "http://example.com/image.jpg"
	store := &entity.Store{
		BaseModel: base_model.BaseModel{CreatedBy: sessionUserID},
		UserID:    sessionUserID,
		StoreName: "TestStore",
		Hours: []*entity.StoreHour{
			{
				Open:  "00:00",
				Close: "00:00",
			},
		},
		Images: []*entity.StoreImage{
			{
				ImageURL:    imgURL,
				ImageBase64: imgBase64,
			},
		},
	}

	store24Hours := &entity.Store{
		BaseModel: base_model.BaseModel{CreatedBy: sessionUserID},
		UserID:    sessionUserID,
		StoreName: "TestStore",
		Hours: []*entity.StoreHour{
			{
				IsOpen: true,
				Is24Hr: true,
				Open:   "00:00",
				Close:  "23:59",
			},
		},
		Images: []*entity.StoreImage{
			{
				ImageURL:    imgURL,
				ImageBase64: imgBase64,
			},
		},
	}

	t.Run("Should return success", func(t *testing.T) {
		storeRepository.EXPECT().
			GetStoreByUserID(ctx, gomock.Any()).
			Times(1).
			Return(nil, nil)

		storage.EXPECT().
			UploadImage(ctx, imgBase64, sessionUserID.String()).
			Return(imgURL, nil)

		storeRepository.EXPECT().
			CreateStore(ctx, gomock.Any()).
			Return(store, nil)

		result, err := service.CreateStore(ctx, store)

		assert.Nil(t, err)
		assert.Equal(t, sessionUserID.String(), result.BaseModel.CreatedBy.String())
		assert.Equal(t, sessionUserID.String(), result.UserID.String())
		assert.Equal(t, store.StoreName, result.StoreName)
		assert.Equal(t, store.Images[0].ImageURL, result.Images[0].ImageURL)
		assert.Equal(t, store.Images[0].ImageBase64, result.Images[0].ImageBase64)
	})

	t.Run("Should return success 24 Hours", func(t *testing.T) {
		storeRepository.EXPECT().
			GetStoreByUserID(ctx, gomock.Any()).
			Times(1).
			Return(nil, nil)

		storage.EXPECT().
			UploadImage(ctx, imgBase64, sessionUserID.String()).
			Return(imgURL, nil)

		storeRepository.EXPECT().
			CreateStore(ctx, gomock.Any()).
			Return(store24Hours, nil)

		result, err := service.CreateStore(ctx, store24Hours)

		assert.Nil(t, err)
		assert.Equal(t, sessionUserID.String(), result.BaseModel.CreatedBy.String())
		assert.Equal(t, sessionUserID.String(), result.UserID.String())
		assert.Equal(t, store24Hours.StoreName, result.StoreName)
		assert.Equal(t, store24Hours.Images[0].ImageURL, result.Images[0].ImageURL)
		assert.Equal(t, store24Hours.Images[0].ImageBase64, result.Images[0].ImageBase64)
	})

	t.Run("Should return error if failed to get store", func(t *testing.T) {
		storeRepository.EXPECT().
			GetStoreByUserID(ctx, gomock.Any()).
			Times(1).
			Return(nil, errors.New("error"))

		_, err := service.CreateStore(ctx, store)

		assert.Error(t, err)
	})

	t.Run("Should return error if failed to upload image", func(t *testing.T) {
		storeRepository.EXPECT().
			GetStoreByUserID(ctx, gomock.Any()).
			Times(1).
			Return(nil, nil)

		storage.EXPECT().
			UploadImage(ctx, imgBase64, sessionUserID.String()).
			Times(1).
			Return("", errors.New("failed to upload image"))

		_, err := service.CreateStore(ctx, store)

		assert.Error(t, err)
	})

	t.Run("Should return error if user has a store", func(t *testing.T) {
		storeRepository.EXPECT().
			GetStoreByUserID(ctx, gomock.Any()).
			Times(1).
			Return(store, nil)

		_, err := service.CreateStore(ctx, store)

		assert.Error(t, err)
	})
}

func TestUpsertProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	mockImageRepo := imageRepoMock.NewMockImageRepository(ctrl)
	ctx := context.Background()

	service := New(mockStoreRepo, mockProdRepo, nil, mockImageRepo)

	storeIdUuid, _ := uuid.Parse(storeID)
	userIdUuid, _ := uuid.Parse(userID)
	otherUserIdUuid, _ := uuid.Parse(otherUserID)
	otherStoreIdUuid, _ := uuid.Parse(otherStoreID)
	productIdUuid, _ := uuid.Parse(productID)
	otherProductID2Uuid, _ := uuid.Parse(otherProductID2)
	roleNames := []string{"merchant"}
	adminRoleNames := []string{"merchant", "admin"}
	var emptyProduct []*prodEntity.Product = nil
	existedProdNames := []string{"bakso aci", "keju"}
	invalidProdNames := []string{"keju", "sepatu"}
	productNames := []string{"indomie", "beras"}

	indomie := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "indomie",
		Uom:           "bungkus",
		ProductTypeID: 1,
		Stock:         1,
	}

	beras := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "beras",
		Uom:           "kg",
		ProductTypeID: 1,
		Stock:         1,
	}

	baksoAci := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "bakso aci",
		Uom:           "pieces",
		ProductTypeID: 1,
		Stock:         0,
	}

	keju := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "keju",
		Uom:           "kg",
		ProductTypeID: 1,
		Stock:         1,
	}

	tas := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "tas",
		Uom:           "",
		ProductTypeID: 2,
		Stock:         1,
	}

	sepatu := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "sepatu",
		Uom:           "pasang",
		ProductTypeID: 2,
		Stock:         1,
	}

	products := []*prodEntity.Product{
		indomie,
		beras,
	}

	updateProduct := []*prodEntity.Product{
		{
			BaseModel: base_model.BaseModel{
				ID: productIdUuid,
			},
			StoreID:       storeIdUuid,
			Name:          "indomie",
			Uom:           "bungkus",
			ProductTypeID: 1,
			Stock:         1,
		},
	}

	productWithoutUOM := []*prodEntity.Product{
		{
			StoreID: storeIdUuid,
			Name:    "keju",
			Stock:   1,
		},
	}

	productWithoutProductType := []*prodEntity.Product{
		{
			StoreID: storeIdUuid,
			Name:    "keju",
			Uom:     "kg",
		},
	}

	productWithNegativeStock := []*prodEntity.Product{
		{
			StoreID:       storeIdUuid,
			Name:          "bantal",
			Uom:           "kg",
			ProductTypeID: 2,
			Stock:         -1,
		},
	}

	updateProductWithoutUOM := []*prodEntity.Product{
		{
			BaseModel: base_model.BaseModel{
				ID: productIdUuid,
			},
			StoreID: storeIdUuid,
			Name:    "keju",
			Stock:   1,
		},
	}

	updateProductWithoutProductType := []*prodEntity.Product{
		{
			BaseModel: base_model.BaseModel{
				ID: productIdUuid,
			},
			StoreID: storeIdUuid,
			Name:    "keju",
			Uom:     "kg",
			Stock:   1,
		},
	}

	updateProductWithNegativeStock := []*prodEntity.Product{
		{
			BaseModel: base_model.BaseModel{
				ID: productIdUuid,
			},
			StoreID:       storeIdUuid,
			Name:          "bantal",
			Uom:           "kg",
			ProductTypeID: 2,
			Stock:         -1,
		},
	}

	updateProductWithInvalidUOM := []*prodEntity.Product{
		{
			BaseModel: base_model.BaseModel{
				ID: productIdUuid,
			},
			StoreID:       storeIdUuid,
			Name:          "keju",
			Uom:           "kg",
			ProductTypeID: 1,
			Stock:         1,
		},
		{
			BaseModel: base_model.BaseModel{
				ID: otherProductID2Uuid,
			},
			StoreID:       storeIdUuid,
			Name:          "tas",
			Uom:           "",
			ProductTypeID: 2,
			Stock:         1,
		},
	}

	prodTypes := []*prodEntity.ProductType{
		{
			BaseMasterDataModel: base_model.BaseMasterDataModel{
				ID: 1,
			},
			Name: "Snack",
		},
	}

	existedProducts := []*prodEntity.Product{
		baksoAci,
		keju,
	}

	invalidUomProduct := []*prodEntity.Product{
		keju,
		tas,
	}

	invalidProdTypeProduct := []*prodEntity.Product{
		keju,
		sepatu,
	}

	updateInvalidProdTypeProduct := []*prodEntity.Product{
		{
			BaseModel: base_model.BaseModel{
				ID: productIdUuid,
			},
			StoreID:       storeIdUuid,
			Name:          "keju",
			Uom:           "kg",
			ProductTypeID: 1,
			Stock:         1,
		},
		{
			BaseModel: base_model.BaseModel{
				ID: otherProductID2Uuid,
			},
			StoreID:       storeIdUuid,
			Name:          "sepatu",
			Uom:           "kg",
			ProductTypeID: 2,
			Stock:         1,
		},
	}

	removeProdImage1 := &prodEntity.ProductImage{
		ProductId:      productIdUuid,
		ImageBase64Str: "aaa",
	}
	removeProdImage2 := &prodEntity.ProductImage{
		ProductId:      productIdUuid,
		ImageBase64Str: "aaa",
	}

	productImagesToBeRemoved := []*prodEntity.ProductImage{
		removeProdImage1,
		removeProdImage2,
	}

	existingProdImagesMap := make(map[uuid.UUID][]*prodEntity.ProductImage)
	existingProdImagesMap[productIdUuid] = append(existingProdImagesMap[productIdUuid], productImagesToBeRemoved...)

	store := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIdUuid,
		},
		UserID: userIdUuid,
	}

	otherStore := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: otherStoreIdUuid,
		},
		UserID: otherUserIdUuid,
	}

	t.Run("Should return error if no product provided", func(t *testing.T) {
		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, emptyProduct...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_NO_PRODUCT_INSERTED.String(), "tidak ada produk yang disimpan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error id not admin", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(&entity.Store{}, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, products...)

		errMsg := util.NewError(codes.PermissionDenied, errPb.StoreErrorCode_DONT_HAVE_PERMISSION_TO_CREATE_OR_UPDATE_STORE.String(), "Anda tidak memiliki izin untuk membuat/memperbarui produk untuk toko ini")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if UOM not provided when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, productWithoutUOM...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_UOM_IS_REQUIRED.String(), "Satuan unit diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product id not provided when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, true, productWithoutUOM...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_IS_REQUIRED.String(), "Product id diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if UOM not provided when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, true, updateProductWithoutUOM...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_UOM_IS_REQUIRED.String(), "Satuan unit diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product type not provided when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, productWithoutProductType...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_TYPE_IS_REQUIRED.String(), "Product type id diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product type not provided when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, true, updateProductWithoutProductType...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_PRODUCT_TYPE_IS_REQUIRED.String(), "Product type id diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if stock negative when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, productWithNegativeStock...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_STOCK_SHOULD_BE_POSITIVE.String(), "Stock harus positif")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if stock negative when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, true, updateProductWithNegativeStock...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_STOCK_SHOULD_BE_POSITIVE.String(), "Stock harus positif")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product already exist when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		mockProdRepo.EXPECT().
			GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, existedProdNames).
			Times(1).
			Return(existedProducts, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, existedProducts...)

		errMsg := util.NewError(codes.AlreadyExists, errPb.StoreErrorCode_PRODUCTS_ARE_ALREADY_REGISTERED.String(), fmt.Sprintf("Produk sudah terdaftar : %s", strings.Join(existedProdNames, ",")))

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if UOM is invalid when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, invalidUomProduct...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_UOM_IS_REQUIRED.String(), "Satuan unit diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if UOM is invalid when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, true, updateProductWithInvalidUOM...)

		errMsg := util.NewError(codes.InvalidArgument, errPb.StoreErrorCode_UOM_IS_REQUIRED.String(), "Satuan unit diperlukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product type invalid when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		mockProdRepo.EXPECT().
			GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, invalidProdNames).
			Times(1).
			Return([]*prodEntity.Product{}, nil)

		mockProdRepo.EXPECT().
			GetProductTypesByIds(gomock.Any(), []int64{1, 2}).
			Times(1).
			Return(prodTypes, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, false, invalidProdTypeProduct...)

		errMsg := util.NewError(codes.NotFound, errPb.StoreErrorCode_PRODUCT_TYPE_ID_IS_NOT_FOUND.String(), "Product type id tidak ditemukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error product type invalid when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), storeID).
			Times(1).
			Return(store, nil)

		mockProdRepo.EXPECT().
			GetProductTypesByIds(gomock.Any(), []int64{1, 2}).
			Times(1).
			Return(prodTypes, nil)

		err := service.UpsertProducts(ctx, userIdUuid, roleNames, storeIdUuid, true, updateInvalidProdTypeProduct...)

		errMsg := util.NewError(codes.NotFound, errPb.StoreErrorCode_PRODUCT_TYPE_ID_IS_NOT_FOUND.String(), "Product type id tidak ditemukan")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return success if different store id but admin when create product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), otherStoreID).
			Times(1).
			Return(otherStore, nil)

		mockProdRepo.EXPECT().
			GetProductsByStoreIdAndNames(gomock.Any(), otherStoreIdUuid, productNames).
			Times(1).
			Return([]*prodEntity.Product{}, nil)

		mockProdRepo.EXPECT().
			GetProductTypesByIds(gomock.Any(), []int64{1}).
			Times(1).
			Return(prodTypes, nil)

		mockProdRepo.EXPECT().
			InitiateTransaction(gomock.Any()).
			Times(1).
			Return(true)

		mockProdRepo.EXPECT().
			UpsertProducts(ctx, gomock.Any()).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			TransactionCommit().
			Times(1).
			Return(nil)

		err := service.UpsertProducts(ctx, userIdUuid, adminRoleNames, otherStoreIdUuid, false, products...)

		assert.Nil(t, err)
	})

	t.Run("Should return success if different store id but admin when update product", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(gomock.Any(), otherStoreID).
			Times(1).
			Return(otherStore, nil)

		mockProdRepo.EXPECT().
			GetProductTypesByIds(gomock.Any(), []int64{1}).
			Times(1).
			Return(prodTypes, nil)

		mockProdRepo.EXPECT().
			InitiateTransaction(gomock.Any()).
			Times(1).
			Return(true)

		mockProdRepo.EXPECT().
			UpsertProducts(ctx, gomock.Any()).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			GetProductImagesByProductIds(ctx, []uuid.UUID{productIdUuid}).
			Times(1).
			Return(productImagesToBeRemoved, existingProdImagesMap, nil)

		mockImageRepo.EXPECT().
			RemoveImage(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			DeleteProductImages(ctx, productImagesToBeRemoved).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			TransactionCommit().
			Times(1).
			Return(nil)

		err := service.UpsertProducts(ctx, userIdUuid, adminRoleNames, otherStoreIdUuid, true, updateProduct...)

		assert.Nil(t, err)
	})
}

func TestUpdateStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	mockImageRepo := imageRepoMock.NewMockImageRepository(ctrl)
	mockStorage := storeRepoMock.NewMockStorage(ctrl)
	service := New(mockStoreRepo, nil, mockStorage, mockImageRepo)

	ctx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	sessionUserID := uuid.New()

	md.Set("x-user-id", sessionUserID.String())
	ctx = metadata.NewIncomingContext(ctx, md)

	storeIDUuid := uuid.MustParse(storeID)
	imgBase64 := "SampleImageBase64"
	imgURL := "http://www.example.com"

	existingStore := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIDUuid,
		},
		UserID:    sessionUserID,
		StoreName: "Toko Maju",
	}

	updatedStore := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIDUuid,
		},
		UserID:    sessionUserID,
		StoreName: "Toko Sebelah",
		Images: []*entity.StoreImage{
			{
				ImageBase64: imgBase64,
				ImageURL:    imgURL,
			},
		},
		Tags: []*entity.StoreTag{
			{
				TagName: "aaa",
			},
		},
		Hours: []*entity.StoreHour{
			{
				Open:  "00:00",
				Close: "00:00",
			},
		},
	}

	updatedStore24Hours := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIDUuid,
		},
		UserID:    sessionUserID,
		StoreName: "Toko Sebelah",
		Images: []*entity.StoreImage{
			{
				ImageBase64: imgBase64,
				ImageURL:    imgURL,
			},
		},
		Tags: []*entity.StoreTag{
			{
				TagName: "aaa",
			},
		},
		Hours: []*entity.StoreHour{
			{
				IsOpen: true,
				Is24Hr: true,
				Open:   "00:00",
				Close:  "23:59",
			},
		},
	}

	t.Run("Should return success", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(ctx, storeID).
			Times(1).
			Return(existingStore, nil)

		mockStorage.EXPECT().
			UploadImage(ctx, imgBase64, storeID).
			Times(1).
			Return(imgURL, nil)

		mockStoreRepo.EXPECT().
			UpdateStore(ctx, gomock.Any()).
			Times(1).
			Return(updatedStore, nil)

		result, err := service.UpdateStore(ctx, storeID, updatedStore)

		assert.Nil(t, err)
		assert.Equal(t, updatedStore, result)
	})

	t.Run("Should return success 24 Hours", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(ctx, storeID).
			Times(1).
			Return(existingStore, nil)

		mockStorage.EXPECT().
			UploadImage(ctx, imgBase64, storeID).
			Times(1).
			Return(imgURL, nil)

		mockStoreRepo.EXPECT().
			UpdateStore(ctx, gomock.Any()).
			Times(1).
			Return(updatedStore24Hours, nil)

		result, err := service.UpdateStore(ctx, storeID, updatedStore24Hours)

		assert.Nil(t, err)
		assert.Equal(t, updatedStore24Hours, result)
	})

	t.Run("Should return error if failed to upload image", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(ctx, storeID).
			Times(1).
			Return(existingStore, nil)

		mockStorage.EXPECT().
			UploadImage(ctx, imgBase64, storeID).
			Times(1).
			Return(imgURL, errors.New("error"))

		updatedStore.Images[0].ImageBase64 = imgBase64
		_, err := service.UpdateStore(ctx, storeID, updatedStore)

		assert.Error(t, err)
	})

	t.Run("Should return error if failed to parse store id", func(t *testing.T) {
		_, err := service.UpdateStore(ctx, "storeID", updatedStore)

		errMsg := status.Errorf(codes.InvalidArgument, "store id should be uuid")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if store not found", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(ctx, storeID).
			Times(1).
			Return(nil, status.Errorf(codes.NotFound, "Store is not found"))

		_, err := service.UpdateStore(ctx, storeID, updatedStore)

		assert.Error(t, err)
	})
}

func TestListStores(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	service := New(mockStoreRepo, nil, nil, nil)

	ctx := context.Background()
	page := 1
	limit := 10

	store := []*entity.Store{
		{
			StoreName: "Toko Sebelah",
		},
	}

	t.Run("Should return success", func(t *testing.T) {
		mockStoreRepo.
			EXPECT().ListStores(ctx, page, limit).
			Times(1).
			Return(store, nil)

		result, err := service.ListStores(ctx, int32(page), int32(limit))

		assert.Nil(t, err)
		assert.Equal(t, store, result)
	})

	t.Run("Should return error if failed to get store", func(t *testing.T) {
		mockStoreRepo.
			EXPECT().ListStores(ctx, page, limit).
			Times(1).
			Return(nil, errors.New("error"))

		_, err := service.ListStores(ctx, int32(page), int32(limit))

		assert.Error(t, err)
	})
}

func TestDeleteStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	service := New(mockStoreRepo, nil, nil, nil)

	ctx := context.Background()

	storeIds := []string{storeID}

	t.Run("Should return success", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			DeleteStores(ctx, storeIds).
			Times(1).
			Return(nil)

		err := service.DeleteStores(ctx, storeIds)

		assert.Nil(t, err)
	})

	t.Run("Should return error if failed to delete store", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			DeleteStores(ctx, storeIds).
			Times(1).
			Return(errors.New("failed"))

		err := service.DeleteStores(ctx, storeIds)

		assert.Error(t, err)
	})
}

func TestUpsertProductCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	service := New(nil, mockProdRepo, nil, nil)

	ctx := context.Background()

	pakaian := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 1,
		},
		Name: "Pakaian",
	}

	kaos := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 2,
		},
		Name: "Pakaian",
	}

	t.Run("Should return error if failed get product category", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryByName(ctx, strings.ToLower(pakaian.Name)).
			Times(1).
			Return(nil, errors.New("failed"))

		err := service.UpsertProductCategory(ctx, pakaian)

		assert.Error(t, err)
	})

	t.Run("Should return error if name already exist", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryByName(ctx, strings.ToLower(pakaian.Name)).
			Times(1).
			Return(kaos, nil)

		err := service.UpsertProductCategory(ctx, pakaian)

		errMsg := status.Errorf(codes.AlreadyExists, "Name is already used by another product category")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if failed upsert product category", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryByName(ctx, strings.ToLower(pakaian.Name)).
			Times(1).
			Return(pakaian, nil)

		mockProdRepo.EXPECT().
			UpsertProductCategory(ctx, pakaian).
			Times(1).
			Return(errors.New("failed"))

		err := service.UpsertProductCategory(ctx, pakaian)

		assert.Error(t, err)
	})

	t.Run("Should return success", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryByName(ctx, strings.ToLower(pakaian.Name)).
			Times(1).
			Return(pakaian, nil)

		mockProdRepo.EXPECT().
			UpsertProductCategory(ctx, pakaian).
			Times(1).
			Return(nil)

		err := service.UpsertProductCategory(ctx, pakaian)

		assert.Nil(t, err)
	})
}

func TestUpsertProductType(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)

	service := New(nil, mockProdRepo, nil, nil)

	ctx := context.Background()

	makanan := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 1,
		},
		Name: "makanan",
	}

	kendaraan := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 2,
		},
		Name: "kendaraan",
	}

	komputer := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 3,
		},
		Name: "komputer",
	}

	indomie := &prodEntity.ProductType{
		Name:              "indomie",
		ProductCategoryID: makanan.BaseMasterDataModel.ID,
	}

	sedan := &prodEntity.ProductType{
		Name:              "sedan",
		ProductCategoryID: kendaraan.BaseMasterDataModel.ID,
	}

	mouse := &prodEntity.ProductType{
		Name:              "mouse",
		ProductCategoryID: komputer.BaseMasterDataModel.ID,
	}

	t.Run("Should return error if name is already exist", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, kendaraan.BaseMasterDataModel.ID).
			Times(1).
			Return(kendaraan, nil)

		mockProdRepo.EXPECT().
			GetProductTypeByName(ctx, sedan.ProductCategoryID, sedan.Name).
			Times(1).
			Return(sedan, nil)

		err := service.UpsertProductType(ctx, sedan)

		errMsg := status.Errorf(codes.AlreadyExists, "Product type is already exist for this product category")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product category not found", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, kendaraan.BaseMasterDataModel.ID).
			Times(1).
			Return(nil, nil)

		err := service.UpsertProductType(ctx, sedan)

		errMsg := status.Errorf(codes.NotFound, "Related product category data is not found")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if failed to update product type", func(t *testing.T) {
		errorMsg := errors.New("error")
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, komputer.BaseMasterDataModel.ID).
			Times(1).
			Return(komputer, nil)

		mockProdRepo.EXPECT().
			GetProductTypeByName(ctx, mouse.ProductCategoryID, mouse.Name).
			Times(1).
			Return(nil, nil)

		mockProdRepo.EXPECT().
			UpsertProductType(ctx, mouse).
			Times(1).
			Return(errorMsg)

		err := service.UpsertProductType(ctx, mouse)

		errMsg := status.Errorf(codes.Internal, "Error when inserting / updating product Type :"+errorMsg.Error())

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return success", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, makanan.BaseMasterDataModel.ID).
			Times(1).
			Return(makanan, nil)

		mockProdRepo.EXPECT().
			GetProductTypeByName(ctx, indomie.ProductCategoryID, indomie.Name).
			Times(1).
			Return(nil, nil)

		mockProdRepo.EXPECT().
			UpsertProductType(ctx, indomie).
			Times(1).
			Return(nil)

		err := service.UpsertProductType(ctx, indomie)

		assert.Nil(t, err)
	})
}

func TestGetProductByStoreId(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)

	service := New(mockStoreRepo, mockProdRepo, nil, nil)

	ctx := context.Background()

	storeIDUuid := uuid.MustParse(storeID)
	otherStoreIDUuid := uuid.MustParse(otherStoreID)
	otherStoreID2Uuid := uuid.MustParse(otherStoreID2)

	store := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIDUuid,
		},
	}

	otherStore2 := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: otherStoreIDUuid,
		},
	}

	products := []*prodEntity.Product{
		{
			Name: "mouse",
		},
	}

	t.Run("Should return error if store id not found", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(ctx, otherStoreID2).
			Times(1).
			Return(nil, status.Errorf(codes.NotFound, "Not Found"))

		_, err := service.GetProductsByStoreId(ctx, otherStoreID2Uuid, nil, false)

		errMsg := status.Errorf(codes.NotFound, "Not Found")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if failed to get product", func(t *testing.T) {
		errorMsg := errors.New("error")
		mockStoreRepo.EXPECT().
			GetStore(ctx, otherStoreID).
			Times(1).
			Return(otherStore2, nil)

		mockProdRepo.EXPECT().
			GetProductsByStoreId(ctx, otherStoreIDUuid, gomock.Any(), false).
			Times(1).
			Return(nil, errorMsg)

		_, err := service.GetProductsByStoreId(ctx, otherStoreIDUuid, nil, false)

		errMsg := status.Errorf(codes.Internal, "Error when getting product list :"+errorMsg.Error())

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return success", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStore(ctx, storeID).
			Times(1).
			Return(store, nil)

		mockProdRepo.EXPECT().
			GetProductsByStoreId(ctx, storeIDUuid, gomock.Any(), false).
			Times(1).
			Return(products, nil)

		result, err := service.GetProductsByStoreId(ctx, storeIDUuid, nil, false)

		assert.Nil(t, err)
		assert.Equal(t, products, result)
	})
}

func TestGetProductCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	service := New(nil, mockProdRepo, nil, nil)

	ctx := context.Background()

	category := []*prodEntity.ProductCategory{
		{
			Name: "shirt",
		},
	}

	t.Run("Should return error if failed to get product categories", func(t *testing.T) {
		errorMsg := errors.New("error")
		mockProdRepo.EXPECT().
			GetProductCategories(ctx, false).
			Return(nil, errorMsg)

		_, _, err := service.GetProductCategories(ctx, false)

		errMsg := status.Errorf(codes.Internal, "Error when getting product categories :"+errorMsg.Error())

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return success if failed to get read file", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategories(ctx, false).
			Return(category, nil)

		cat, uom, err := service.GetProductCategories(ctx, false)

		assert.Nil(t, err)
		assert.Nil(t, uom)
		assert.Equal(t, category, cat)
	})

	t.Run("Should return success", func(t *testing.T) {
		// write file
		uoms := []string{"pieces", "kilogram", "ons", "pound", "botol"}
		fileName := "uom.json"
		file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer func() {
			file.Close()
			_ = os.Remove("uom.json")
		}()
		encoder := json.NewEncoder(file)
		_ = encoder.Encode(uoms)

		mockProdRepo.EXPECT().
			GetProductCategories(ctx, false).
			Return(category, nil)

		cat, uom, err := service.GetProductCategories(ctx, false)

		assert.Nil(t, err)
		assert.Equal(t, uoms, uom)
		assert.Equal(t, category, cat)
	})
}

func TestGetProductTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)

	service := New(nil, mockProdRepo, nil, nil)

	ctx := context.Background()

	prodCat1 := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 1,
		},
	}

	prodCat2 := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 2,
		},
	}

	prodTypes := []*prodEntity.ProductType{
		{
			Name: "makanan",
		},
	}

	t.Run("Should return error if failed to get product category", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, prodCat2.ID).
			Times(1).
			Return(nil, errors.New("error"))

		_, err := service.GetProductTypes(ctx, prodCat2.ID, false)

		errMsg := status.Errorf(codes.AlreadyExists, "Error getting product category by id data")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product category not exist", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, prodCat2.ID).
			Times(1).
			Return(nil, nil)

		_, err := service.GetProductTypes(ctx, prodCat2.ID, false)

		errMsg := status.Errorf(codes.NotFound, "Product category id is not found")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if failed to get product types", func(t *testing.T) {
		errorMsg := errors.New("error")
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, prodCat1.ID).
			Times(1).
			Return(prodCat1, nil)

		mockProdRepo.EXPECT().
			GetProductTypes(ctx, prodCat1.ID, false).
			Return(nil, errorMsg)

		_, err := service.GetProductTypes(ctx, prodCat1.ID, false)

		errMsg := status.Errorf(codes.Internal, "Error when getting product types :"+errorMsg.Error())

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return success", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductCategoryById(ctx, prodCat1.ID).
			Times(1).
			Return(prodCat1, nil)

		mockProdRepo.EXPECT().
			GetProductTypes(ctx, prodCat1.ID, false).
			Return(prodTypes, nil)

		result, err := service.GetProductTypes(ctx, prodCat1.ID, false)

		assert.Nil(t, err)
		assert.Equal(t, prodTypes, result)
	})
}

func TestGetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	service := New(nil, mockProdRepo, nil, nil)

	ctx := context.Background()

	productIDUuid := uuid.MustParse(productID)
	otherProductIDUuid := uuid.MustParse(otherProductID)
	otherProductIDUuid2 := uuid.MustParse(otherProductID2)

	products := &prodEntity.Product{
		Name: "mouse",
	}

	t.Run("Should return error if failed to get product by id", func(t *testing.T) {
		errorMsg := errors.New("error")
		mockProdRepo.EXPECT().
			GetProductById(ctx, otherProductIDUuid).
			Times(1).
			Return(nil, errorMsg)

		_, err := service.GetProductById(ctx, otherProductIDUuid)

		errMsg := status.Errorf(codes.Internal, "Error when getting product by id :"+errorMsg.Error())

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return error if product id not found", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductById(ctx, otherProductIDUuid2).
			Times(1).
			Return(nil, nil)

		_, err := service.GetProductById(ctx, otherProductIDUuid2)

		errMsg := status.Errorf(codes.NotFound, "Product id not found")

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})

	t.Run("Should return success", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductById(ctx, productIDUuid).
			Times(1).
			Return(products, nil)

		result, err := service.GetProductById(ctx, productIDUuid)

		assert.Nil(t, err)
		assert.Equal(t, products, result)
	})
}

func TestDeleteProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	mockImageRepo := imageRepoMock.NewMockImageRepository(ctrl)

	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)

	productService := New(nil, mockProdRepo, nil, mockImageRepo)

	productIDUuid := uuid.MustParse(productID)
	userIDUuid := uuid.MustParse(userID)
	removeProdImage1 := &prodEntity.ProductImage{
		ProductId:      productIDUuid,
		ImageBase64Str: "aaa",
	}
	removeProdImage2 := &prodEntity.ProductImage{
		ProductId:      productIDUuid,
		ImageBase64Str: "aaa",
	}

	productImagesToBeRemoved := []*prodEntity.ProductImage{removeProdImage1, removeProdImage2}
	existingProdImagesMap := make(map[uuid.UUID][]*prodEntity.ProductImage)

	t.Run("Should return empty when success", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductImagesByProductIds(ctx, []uuid.UUID{productIDUuid}).
			Times(1).
			Return(productImagesToBeRemoved, existingProdImagesMap, nil)

		mockProdRepo.EXPECT().
			InitiateTransaction(gomock.Any()).
			Times(1).
			Return(true)

		mockImageRepo.EXPECT().
			RemoveImage(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			DeleteProductImages(ctx, productImagesToBeRemoved).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			DeleteProductById(gomock.Any(), productIDUuid).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			TransactionCommit().
			Times(1)

		err := productService.DeleteProductById(ctx, userIDUuid, productIDUuid)

		assert.NoError(t, err)
	})

	t.Run("Should return error when failed to get product image", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductImagesByProductIds(ctx, []uuid.UUID{productIDUuid}).
			Times(1).
			Return([]*prodEntity.ProductImage{}, existingProdImagesMap, err)

		err := productService.DeleteProductById(ctx, userIDUuid, productIDUuid)

		assert.Error(t, err)
	})

	t.Run("Should return error when failed delete image", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductImagesByProductIds(ctx, []uuid.UUID{productIDUuid}).
			Times(1).
			Return(productImagesToBeRemoved, existingProdImagesMap, nil)

		mockProdRepo.EXPECT().
			InitiateTransaction(gomock.Any()).
			Times(1).
			Return(true)

		mockImageRepo.EXPECT().
			RemoveImage(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			DeleteProductImages(ctx, productImagesToBeRemoved).
			Times(1).
			Return(err)

		mockProdRepo.EXPECT().
			TransactionRollback().
			Times(1)

		err := productService.DeleteProductById(ctx, userIDUuid, productIDUuid)

		assert.Error(t, err)
	})

	t.Run("Should return error when failed remove image", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductImagesByProductIds(ctx, []uuid.UUID{productIDUuid}).
			Times(1).
			Return(productImagesToBeRemoved, existingProdImagesMap, nil)

		mockProdRepo.EXPECT().
			InitiateTransaction(gomock.Any()).
			Times(1).
			Return(true)

		mockImageRepo.EXPECT().
			RemoveImage(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(err)

		mockProdRepo.EXPECT().
			TransactionRollback().
			Times(1)

		err := productService.DeleteProductById(ctx, userIDUuid, productIDUuid)

		assert.Error(t, err)
	})

	t.Run("Should return error when failed to delete", func(t *testing.T) {
		mockProdRepo.EXPECT().
			GetProductImagesByProductIds(ctx, []uuid.UUID{productIDUuid}).
			Times(1).
			Return(productImagesToBeRemoved, existingProdImagesMap, nil)

		mockProdRepo.EXPECT().
			InitiateTransaction(gomock.Any()).
			Times(1).
			Return(true)

		mockImageRepo.EXPECT().
			RemoveImage(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			DeleteProductImages(ctx, productImagesToBeRemoved).
			Times(1).
			Return(nil)

		mockProdRepo.EXPECT().
			DeleteProductById(gomock.Any(), productIDUuid).
			Times(1).
			Return(err)

		mockProdRepo.EXPECT().
			TransactionRollback().
			Times(1)

		err := productService.DeleteProductById(ctx, userIDUuid, productIDUuid)

		assert.Error(t, err)
	})
}

func TestGetStoreByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	service := New(mockStoreRepo, nil, nil, nil)

	ctx := context.Background()

	storeIdUuid, _ := uuid.Parse(storeID)
	store := entity.Store{
		StoreName: "Toko Agak Laen",
	}

	t.Run("Should return success", func(t *testing.T) {
		mockStoreRepo.EXPECT().
			GetStoreByUserID(ctx, storeIdUuid).
			Times(1).
			Return(&store, nil)

		result, err := service.GetStoreByUserID(ctx, storeIdUuid)

		assert.Nil(t, err)
		assert.Equal(t, &store, result)
	})

	t.Run("Should return error if failed to get store by user id", func(t *testing.T) {
		errorMsg := errors.New("error")
		mockStoreRepo.EXPECT().
			GetStoreByUserID(ctx, storeIdUuid).
			Times(1).
			Return(nil, errorMsg)

		_, err := service.GetStoreByUserID(ctx, storeIdUuid)

		errMsg := status.Errorf(codes.Internal, "Error when getting store by user id :"+errorMsg.Error())

		assert.Error(t, err)
		assert.Equal(t, errMsg, err)
	})
}
