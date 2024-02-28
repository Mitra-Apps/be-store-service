package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	prodRepoMock "github.com/Mitra-Apps/be-store-service/domain/product/repository/mock"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	storeRepoMock "github.com/Mitra-Apps/be-store-service/domain/store/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	userID          = "8b15140c-f6d0-4f2f-8302-57383a51adaf"
	otherUserID     = "2f27d467-9f83-4170-96ab-36e0994f37ca"
	storeID         = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
	otherStoreID    = "52d11042-8c45-453e-86af-fe1e4d7facf6"
	otherStoreID2   = "52d11042-8c45-453e-86af-fe1e4d7facf7"
	productID       = "7d56be32-70a2-4f49-b66b-63e6f8e719d7"
	otherProductID  = "7d56be32-70a2-4f49-b66b-63e6f8e719d8"
	otherProductID2 = "7d56be32-70a2-4f49-b66b-63e6f8e719d9"
)

func Test_service_OpenCloseStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	ctx := context.Background()
	userIdUuid, _ := uuid.Parse(userID)
	otherUserIdUuid, _ := uuid.Parse(otherUserID)
	storeIdUuid, _ := uuid.Parse(storeID)
	roleNames := []string{}
	roleNames = append(roleNames, "merchant")
	adminRole := []string{}
	adminRole = append(adminRole, "admin")
	mockStoreRepo.EXPECT().GetStore(gomock.Any(), gomock.Any()).Return(&entity.Store{
		UserID: userIdUuid,
	}, nil).AnyTimes()
	mockStoreRepo.EXPECT().OpenCloseStore(ctx, storeIdUuid, false).Return(nil).AnyTimes()

	type fields struct {
		storeRepository *storeRepoMock.MockStoreServiceRepository
	}
	type args struct {
		ctx       context.Context
		userID    uuid.UUID
		roleNames []string
		storeID   string
		isActive  bool
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "OpenCloseStore_NoStoreIdProvided_ReturnValidationError",
			fields: fields{
				storeRepository: mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				storeID:   "",
				isActive:  false,
				roleNames: roleNames,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.InvalidArgument, "store id is required"),
		},
		{
			name: "OpenCloseStore_StoreIdIsNotUUID_ReturnValidationError",
			fields: fields{
				storeRepository: mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   "aaa",
				isActive:  false,
				roleNames: roleNames,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.InvalidArgument, "store id should be uuid"),
		},
		{
			name: "OpenCloseStore_DifferentUserIDNotAdmin_DontHavePermission",
			fields: fields{
				storeRepository: mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    otherUserIdUuid,
				storeID:   storeID,
				isActive:  false,
				roleNames: roleNames,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.PermissionDenied, "You do not have permission to open / close this store"),
		},
		{
			name: "OpenCloseStore_DifferentUserIDAdmin_DontHavePermission",
			fields: fields{
				storeRepository: mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    otherUserIdUuid,
				storeID:   storeID,
				isActive:  false,
				roleNames: adminRole,
			},
			wantErr: false,
		},
		{
			name: "OpenCloseStore_NoError_Success",
			fields: fields{
				storeRepository: mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   storeID,
				isActive:  false,
				roleNames: roleNames,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.storeRepository, nil, nil, nil)
			if err := s.OpenCloseStore(tt.args.ctx, tt.args.userID, tt.args.roleNames, tt.args.storeID, tt.args.isActive); tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCreateStore(t *testing.T) {
	ctx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	sessionUserID := uuid.New()

	md.Set("x-user-id", sessionUserID.String())
	ctx = metadata.NewIncomingContext(ctx, md)

	testCases := []struct {
		name          string
		setupMocks    func(storeRepository *storeRepoMock.MockStoreServiceRepository, storage *storeRepoMock.MockStorage)
		inputStore    *entity.Store
		expectedStore *entity.Store
		expectedError error
	}{
		{
			name: "Successful store creation",
			setupMocks: func(storeRepository *storeRepoMock.MockStoreServiceRepository, storage *storeRepoMock.MockStorage) {
				storeRepository.EXPECT().GetStoreByUserID(ctx, gomock.Any()).Return(nil, nil)
				storage.EXPECT().UploadImage(ctx, "SampleImageBase64", sessionUserID.String()).Return("http://example.com/image.jpg", nil)
				storeRepository.EXPECT().CreateStore(ctx, gomock.Any()).Return(&entity.Store{
					BaseModel: base_model.BaseModel{CreatedBy: sessionUserID},
					UserID:    sessionUserID,
					StoreName: "TestStore",
					Images: []*entity.StoreImage{
						{
							ImageURL:    "http://example.com/image.jpg",
							ImageBase64: "",
						},
					},
				}, nil)
			},
			inputStore: &entity.Store{
				UserID:    sessionUserID,
				StoreName: "TestStore",
				Images: []*entity.StoreImage{
					{
						ImageBase64: "SampleImageBase64",
					},
				},
			},
			expectedStore: &entity.Store{
				BaseModel: base_model.BaseModel{CreatedBy: sessionUserID},
				UserID:    sessionUserID,
				StoreName: "TestStore",
				Images: []*entity.StoreImage{
					{
						ImageURL:    "http://example.com/image.jpg",
						ImageBase64: "",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Error uploading image",
			setupMocks: func(storeRepository *storeRepoMock.MockStoreServiceRepository, storage *storeRepoMock.MockStorage) {
				storeRepository.EXPECT().GetStoreByUserID(ctx, gomock.Any()).Return(nil, nil)
				storage.EXPECT().
					UploadImage(ctx, "SampleImageBase64_failed", sessionUserID.String()).
					Return("", errors.New("failed to upload image"))
			},
			inputStore: &entity.Store{
				StoreName: "TestStore",
				Images: []*entity.StoreImage{
					{
						ImageBase64: "SampleImageBase64_failed",
					},
				},
			},
			expectedStore: nil,
			expectedError: errors.New("failed to upload image"),
		},
		{
			name: "User already has store",
			setupMocks: func(storeRepository *storeRepoMock.MockStoreServiceRepository, storage *storeRepoMock.MockStorage) {
				storeRepository.EXPECT().GetStoreByUserID(ctx, gomock.Any()).Return(&entity.Store{
					StoreName: "TestStore",
				}, nil)
			},
			inputStore: &entity.Store{
				StoreName: "TestStore",
			},
			expectedStore: nil,
			expectedError: status.Errorf(codes.AlreadyExists, "User already has a store"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storeRepository := storeRepoMock.NewMockStoreServiceRepository(ctrl)
			storage := storeRepoMock.NewMockStorage(ctrl)
			service := New(storeRepository, nil, storage, nil)

			tc.setupMocks(storeRepository, storage)
			resultStore, err := service.CreateStore(ctx, tc.inputStore)
			assert.Equal(t, tc.expectedStore, resultStore)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func Test_service_UpsertProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	ctx := context.Background()
	userIdUuid, _ := uuid.Parse(userID)
	otherUserIdUuid, _ := uuid.Parse(otherUserID)
	storeIdUuid, _ := uuid.Parse(storeID)
	otherStoreIdUuid, _ := uuid.Parse(otherStoreID)
	products := []*prodEntity.Product{}
	indomie := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "indomie",
		UomID:         1,
		ProductTypeID: 1,
	}
	beras := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "beras",
		UomID:         1,
		ProductTypeID: 1,
	}
	baksoAci := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "bakso aci",
		UomID:         1,
		ProductTypeID: 1,
	}
	keju := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "keju",
		UomID:         1,
		ProductTypeID: 1,
	}
	tas := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "tas",
		UomID:         2,
		ProductTypeID: 2,
	}
	sepatu := &prodEntity.Product{
		StoreID:       storeIdUuid,
		Name:          "sepatu",
		UomID:         1,
		ProductTypeID: 2,
	}
	products = append(products, indomie)
	products = append(products, beras)
	existedProducts := []*prodEntity.Product{}
	existedProducts = append(existedProducts, baksoAci)
	existedProducts = append(existedProducts, keju)
	invalidUomProduct := []*prodEntity.Product{
		keju,
		tas,
	}
	invalidProdTypeProduct := []*prodEntity.Product{
		keju,
		sepatu,
	}
	noUOM := []*prodEntity.Product{
		&prodEntity.Product{
			StoreID: storeIdUuid,
			Name:    "keju",
		},
	}
	noProdType := []*prodEntity.Product{
		&prodEntity.Product{
			StoreID: storeIdUuid,
			Name:    "keju",
			UomID:   1,
		},
	}
	roleNames := []string{"merchant"}
	adminRoleNames := []string{"merchant", "admin"}
	existedProdNames := []string{"bakso aci", "keju"}
	uoms := []*prodEntity.UnitOfMeasure{
		&prodEntity.UnitOfMeasure{
			BaseMasterDataModel: base_model.BaseMasterDataModel{
				ID: 1,
			},
			Name: "kilogram",
		},
	}
	prodTypes := []*prodEntity.ProductType{
		&prodEntity.ProductType{
			BaseMasterDataModel: base_model.BaseMasterDataModel{
				ID: 1,
			},
			Name: "Snack",
		},
	}
	mockProdRepo.EXPECT().GetUnitOfMeasuresByIds(gomock.Any(), []int64{1}).Return(uoms, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasuresByIds(gomock.Any(), []int64{1, 2}).Return(uoms, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypesByIds(gomock.Any(), []int64{1}).Return(prodTypes, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypesByIds(gomock.Any(), []int64{1, 2}).Return(prodTypes, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, existedProdNames).Return(existedProducts, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, []string{"keju", "tas"}).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, []string{"keju", "sepatu"}).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, []string{"indomie", "beras"}).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreIdAndNames(gomock.Any(), otherStoreIdUuid, []string{"indomie", "beras"}).Return(nil, nil).AnyTimes()

	mockProdRepo.EXPECT().UpsertProducts(ctx, products).Return(nil).AnyTimes()
	mockStoreRepo.EXPECT().GetStore(gomock.Any(), otherStoreID).Return(&entity.Store{
		BaseModel: base_model.BaseModel{
			ID: otherStoreIdUuid,
		},
		UserID: otherUserIdUuid,
	}, nil).AnyTimes()
	mockStoreRepo.EXPECT().GetStore(gomock.Any(), storeID).Return(&entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIdUuid,
		},
		UserID: userIdUuid,
	}, nil).AnyTimes()

	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
		storeRepository   *storeRepoMock.MockStoreServiceRepository
	}
	type args struct {
		ctx       context.Context
		userID    uuid.UUID
		roleNames []string
		storeID   uuid.UUID
		products  []*prodEntity.Product
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "UpsertProduct_NoProductProvided_ReturnValidationError",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  nil,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.InvalidArgument, "No product inserted"),
		},
		{
			name: "UpsertProduct_DifferenStoreIDNotAdmin_DontHavePermission",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   otherStoreIdUuid,
				roleNames: roleNames,
				products:  products,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.PermissionDenied, "You don't have permission to create / update product for this store"),
		},
		{
			name: "UpsertProduct_UomNotProvided_Error",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  noUOM,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.InvalidArgument, "Uom id is required"),
		},
		{
			name: "UpsertProduct_ProdTypeNotProvided_Error",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  noProdType,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.InvalidArgument, "Product type id is required"),
		},
		{
			name: "UpsertProduct_ProductAlreadyExisted_ReturnValidationError",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  existedProducts,
				userID:    userIdUuid,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.AlreadyExists, "Product are already exist : "+strings.Join(existedProdNames, ",")),
		},
		{
			name: "UpsertProduct_InvalidUOM_Error",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  invalidUomProduct,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.NotFound, "Unit of measure id is not found"),
		},
		{
			name: "UpsertProduct_InvalidProductType_Error",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  invalidProdTypeProduct,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.NotFound, "Product type id is not found"),
		},
		{
			name: "UpsertProduct_DifferenStoreIDButAdmin_Success",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   otherStoreIdUuid,
				roleNames: adminRoleNames,
				products:  products,
			},
			wantErr: false,
		},
		{
			name: "UpsertProduct_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:       ctx,
				userID:    userIdUuid,
				storeID:   storeIdUuid,
				roleNames: roleNames,
				products:  products,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.storeRepository, tt.fields.productRepository, nil, nil)
			if err := s.UpsertProducts(tt.args.ctx, tt.args.userID, tt.args.roleNames, tt.args.storeID, tt.args.products); tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUpdateStore(t *testing.T) {
	ctx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	sessionUserID := uuid.New()

	md.Set("x-user-id", sessionUserID.String())
	ctx = metadata.NewIncomingContext(ctx, md)

	testCases := []struct {
		name       string
		setupMocks func(storeRepository *storeRepoMock.MockStoreServiceRepository, storage *storeRepoMock.MockStorage)
		inputStore struct {
			storeID string
			store   *entity.Store
		}
		expectedStore *entity.Store
		expectedError error
	}{
		{
			name: "Success",
			setupMocks: func(storeRepository *storeRepoMock.MockStoreServiceRepository, storage *storeRepoMock.MockStorage) {
				storeRepository.EXPECT().UpdateStore(ctx, gomock.Any()).Return(nil)
			},
			inputStore: struct {
				storeID string
				store   *entity.Store
			}{
				storeID: "TestStore",
				store: &entity.Store{
					StoreName: "TestStore",
				},
			},
			expectedStore: &entity.Store{
				StoreName: "TestStore",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storeRepository := storeRepoMock.NewMockStoreServiceRepository(ctrl)
			storage := storeRepoMock.NewMockStorage(ctrl)
			service := New(storeRepository, nil, storage, nil)

			tc.setupMocks(storeRepository, storage)
			result, err := service.UpdateStore(ctx, tc.inputStore.storeID, tc.inputStore.store)
			assert.Equal(t, tc.expectedError, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func Test_service_UpsertUnitOfMeasure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	pcs := &prodEntity.UnitOfMeasure{
		Name:   "pieces",
		Symbol: "pcs",
	}
	gram := &prodEntity.UnitOfMeasure{
		Name:   "gram",
		Symbol: "g",
	}
	kg := &prodEntity.UnitOfMeasure{
		Name:   "kilogram",
		Symbol: "kg",
	}
	ton := &prodEntity.UnitOfMeasure{
		Name:   "ton",
		Symbol: "ton",
	}
	errMsg := "ERROR"
	err := errors.New(errMsg)
	mockProdRepo.EXPECT().GetUnitOfMeasureByName(ctx, pcs.Name).Return(pcs, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasureByName(ctx, gram.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasureByName(ctx, kg.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasureByName(ctx, ton.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasureBySymbol(ctx, gram.Symbol).Return(gram, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasureBySymbol(ctx, kg.Symbol).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasureBySymbol(ctx, ton.Symbol).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().UpsertUnitOfMeasure(ctx, kg).Return(err).AnyTimes()
	mockProdRepo.EXPECT().UpsertUnitOfMeasure(ctx, ton).Return(nil).AnyTimes()

	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx context.Context
		uom *prodEntity.UnitOfMeasure
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "UpsertUnitOfMeasure_NameIsExist_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx: ctx,
				uom: pcs,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.AlreadyExists, "Uom name is already exist in database"),
		},
		{
			name: "UpsertUnitOfMeasure_SymbolIsExist_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx: ctx,
				uom: gram,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.AlreadyExists, "Uom symbol is already exist in database"),
		},
		{
			name: "UpsertUnitOfMeasure_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx: ctx,
				uom: kg,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when inserting / updating unit of measure :"+errMsg),
		},
		{
			name: "UpsertUnitOfMeasure_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx: ctx,
				uom: ton,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if err := s.UpsertUnitOfMeasure(tt.args.ctx, tt.args.uom); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_service_UpsertProductCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	pakaian := &prodEntity.ProductCategory{
		Name: "Pakaian",
	}
	komputer := &prodEntity.ProductCategory{
		Name: "Komputer",
	}
	makanan := &prodEntity.ProductCategory{
		Name: "Makanan",
	}
	errMsg := "ERROR"
	err := errors.New(errMsg)
	mockProdRepo.EXPECT().GetProductCategoryByName(ctx, pakaian.Name).Return(pakaian, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductCategoryByName(ctx, komputer.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductCategoryByName(ctx, makanan.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().UpsertProductCategory(ctx, komputer).Return(err).AnyTimes()
	mockProdRepo.EXPECT().UpsertProductCategory(ctx, makanan).Return(nil).AnyTimes()

	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx             context.Context
		productCategory *prodEntity.ProductCategory
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "UpsertProductCategory_NameAlreadyExist_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:             ctx,
				productCategory: pakaian,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.AlreadyExists, "Category name is already exist in database"),
		},
		{
			name: "UpsertProductCategory_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:             ctx,
				productCategory: komputer,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when inserting / updating product category :"+errMsg),
		},
		{
			name: "UpsertProductCategory_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:             ctx,
				productCategory: makanan,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if err := s.UpsertProductCategory(tt.args.ctx, tt.args.productCategory); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_service_UpsertProductType(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	makanan := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 1,
		},
		Name: "makanan",
	}
	komputer := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 3,
		},
		Name: "komputer",
	}
	kendaraan := &prodEntity.ProductCategory{
		BaseMasterDataModel: base_model.BaseMasterDataModel{
			ID: 2,
		},
		Name: "kendaraan",
	}
	mouse := &prodEntity.ProductType{
		Name:              "mouse",
		ProductCategoryID: komputer.BaseMasterDataModel.ID,
	}
	indomie := &prodEntity.ProductType{
		Name:              "indomie",
		ProductCategoryID: makanan.BaseMasterDataModel.ID,
	}
	pizza := &prodEntity.ProductType{
		Name:              "pizza",
		ProductCategoryID: makanan.BaseMasterDataModel.ID,
	}
	sedan := &prodEntity.ProductType{
		Name:              "sedan",
		ProductCategoryID: kendaraan.BaseMasterDataModel.ID,
	}
	errMsg := "ERROR"
	err := errors.New(errMsg)
	mockProdRepo.EXPECT().GetProductCategoryById(ctx, makanan.BaseMasterDataModel.ID).Return(makanan, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductCategoryById(ctx, komputer.BaseMasterDataModel.ID).Return(komputer, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductCategoryById(ctx, kendaraan.BaseMasterDataModel.ID).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypeByName(ctx, gomock.Any(), pizza.Name).Return(pizza, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypeByName(ctx, gomock.Any(), mouse.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypeByName(ctx, gomock.Any(), indomie.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypeByName(ctx, gomock.Any(), kendaraan.Name).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().UpsertProductType(ctx, mouse).Return(err).AnyTimes()
	mockProdRepo.EXPECT().UpsertProductType(ctx, indomie).Return(nil).AnyTimes()

	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx         context.Context
		productType *prodEntity.ProductType
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "UpsertProductType_NameIsAlreadyExist_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:         ctx,
				productType: sedan,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.AlreadyExists, "Product type is already exist for this product category"),
		},
		{
			name: "UpsertProductType_ProductCategoryNotFound_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:         ctx,
				productType: sedan,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.NotFound, "Related product category data is not found"),
		},
		{
			name: "UpsertProductType_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:         ctx,
				productType: mouse,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when inserting / updating product Type :"+errMsg),
		},
		{
			name: "UpsertProductType_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:         ctx,
				productType: indomie,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if err := s.UpsertProductType(tt.args.ctx, tt.args.productType); err != nil && tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_service_GetUnitOfMeasures(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)
	uoms := []*prodEntity.UnitOfMeasure{}
	uoms = append(uoms, &prodEntity.UnitOfMeasure{
		Name: "kg",
	})
	mockProdRepo.EXPECT().GetUnitOfMeasures(ctx, false).Return(nil, err).AnyTimes()
	mockProdRepo.EXPECT().GetUnitOfMeasures(ctx, true).Return(uoms, nil).AnyTimes()
	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx                  context.Context
		isIncludeDeactivated bool
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "GetUnitOfMeasures_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				isIncludeDeactivated: false,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when getting unit of measures :"+errMsg),
		},
		{
			name: "GetUnitOfMeasures_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				isIncludeDeactivated: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if uom, err := s.GetUnitOfMeasures(tt.args.ctx, tt.args.isIncludeDeactivated); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, uom)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, uom)
			}
		})
	}
}

func Test_service_GetProductsByStoreId(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)
	storeIDUuid := uuid.MustParse(storeID)
	otherStoreIDUuid := uuid.MustParse(otherStoreID)
	otherStoreID2Uuid := uuid.MustParse(otherStoreID2)

	products := []*prodEntity.Product{
		&prodEntity.Product{
			Name: "mouse",
		},
	}

	mockStoreRepo.EXPECT().GetStore(ctx, storeID).Return(&entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeIDUuid,
		},
	}, nil).AnyTimes()
	mockStoreRepo.EXPECT().GetStore(ctx, otherStoreID).Return(&entity.Store{
		BaseModel: base_model.BaseModel{
			ID: otherStoreIDUuid,
		},
	}, nil).AnyTimes()
	mockStoreRepo.EXPECT().GetStore(ctx, otherStoreID2).Return(nil, status.Errorf(codes.NotFound, "Not Found")).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreId(ctx, otherStoreIDUuid, gomock.Any(), false).Return(nil, err).AnyTimes()
	mockProdRepo.EXPECT().GetProductsByStoreId(ctx, storeIDUuid, gomock.Any(), false).Return(products, nil).AnyTimes()
	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
		storeRepository   *storeRepoMock.MockStoreServiceRepository
	}
	type args struct {
		ctx                  context.Context
		storeID              uuid.UUID
		productTypeId        *int64
		isIncludeDeactivated bool
	}
	tests := []struct {
		name          string
		s             *service
		fields        fields
		args          args
		expectedError error
		wantErr       bool
	}{
		{
			name: "GetProductsByStoreId_StoreIDNotFound_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:     ctx,
				storeID: otherStoreID2Uuid,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.NotFound, "Not Found"),
		},
		{
			name: "GetProductsByStoreId_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:     ctx,
				storeID: otherStoreIDUuid,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when getting product list :"+errMsg),
		},
		{
			name: "GetProductsByStoreId_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
				storeRepository:   mockStoreRepo,
			},
			args: args{
				ctx:     ctx,
				storeID: storeIDUuid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.storeRepository, tt.fields.productRepository, nil, nil)
			if gotProducts, err := s.GetProductsByStoreId(tt.args.ctx, tt.args.storeID, tt.args.productTypeId, tt.args.isIncludeDeactivated); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, gotProducts)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, gotProducts)
			}
		})
	}
}

func Test_service_GetProductCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)
	uoms := []*prodEntity.ProductCategory{}
	uoms = append(uoms, &prodEntity.ProductCategory{
		Name: "kg",
	})
	mockProdRepo.EXPECT().GetProductCategories(ctx, false).Return(nil, err).AnyTimes()
	mockProdRepo.EXPECT().GetProductCategories(ctx, true).Return(uoms, nil).AnyTimes()
	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx                  context.Context
		isIncludeDeactivated bool
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "GetProductCategories_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				isIncludeDeactivated: false,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when getting unit of measures :"+errMsg),
		},
		{
			name: "GetProductCategories_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				isIncludeDeactivated: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if uom, err := s.GetProductCategories(tt.args.ctx, tt.args.isIncludeDeactivated); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, uom)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, uom)
			}
		})
	}
}

func Test_service_GetProductTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)
	prodTypes := []*prodEntity.ProductType{}
	prodTypes = append(prodTypes, &prodEntity.ProductType{
		Name: "makanan",
	})
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

	mockProdRepo.EXPECT().GetProductCategoryById(ctx, prodCat2.ID).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductCategoryById(ctx, prodCat1.ID).Return(prodCat1, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypes(ctx, prodCat1.ID, false).Return(nil, err).AnyTimes()
	mockProdRepo.EXPECT().GetProductTypes(ctx, prodCat1.ID, true).Return(prodTypes, nil).AnyTimes()
	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx                  context.Context
		productCategoryId    int64
		isIncludeDeactivated bool
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "GetProductTypes_ProdCategoryNotExist_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				productCategoryId:    prodCat2.ID,
				isIncludeDeactivated: false,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.NotFound, "Product category id is not found"),
		},
		{
			name: "GetProductTypes_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				productCategoryId:    prodCat1.ID,
				isIncludeDeactivated: false,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when getting product types :"+errMsg),
		},
		{
			name: "GetProductTypes_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:                  ctx,
				productCategoryId:    prodCat1.ID,
				isIncludeDeactivated: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if prodType, err := s.GetProductTypes(tt.args.ctx, tt.args.productCategoryId, tt.args.isIncludeDeactivated); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, prodType)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, prodType)
			}
		})
	}
}

func Test_service_UpdateUnitOfMeasure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)
	var uomId1 int64 = 1
	var uomId2 int64 = 2
	var uomId3 int64 = 3

	initialUom := &prodEntity.UnitOfMeasure{
		Name:     "celcius",
		Symbol:   "c",
		IsActive: true,
	}

	updatedUom1 := &prodEntity.UnitOfMeasure{
		Name:     "fahrenheit",
		Symbol:   "f",
		IsActive: true,
	}

	updatedUom2 := &prodEntity.UnitOfMeasure{
		Name:     "kelvin",
		Symbol:   "k",
		IsActive: true,
	}

	// first scenario: ERROR getting UoM by ID
	mockProdRepo.EXPECT().GetUnitOfMeasureById(ctx, uomId1).Return(nil, err)

	// second scenario: ERROR updating the UoM
	mockProdRepo.EXPECT().GetUnitOfMeasureById(ctx, uomId2).Return(initialUom, nil)
	mockProdRepo.EXPECT().UpsertUnitOfMeasure(ctx, updatedUom1).Return(err).AnyTimes()

	// third scenario: SUCCESS updating the UoM
	mockProdRepo.EXPECT().GetUnitOfMeasureById(ctx, uomId3).Return(initialUom, nil)
	mockProdRepo.EXPECT().UpsertUnitOfMeasure(ctx, updatedUom2).Return(nil).AnyTimes()

	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}

	type args struct {
		ctx   context.Context
		uomId int64
		uom   *prodEntity.UnitOfMeasure
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "UpdateUnitOfMeasure_Error_UnableToGetUoMByID",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:   ctx,
				uomId: uomId1,
				uom:   updatedUom1,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when getting uom: "+errMsg),
		},
		{
			name: "UpdateUnitOfMeasure_Error_UnableToUpdateUoM",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:   ctx,
				uomId: uomId2,
				uom:   updatedUom1,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when updating unit of measure: "+errMsg),
		},
		{
			name: "UpdateUnitOfMeasure_NoError_SuccessUpdatingUoM",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:   ctx,
				uomId: uomId3,
				uom:   updatedUom2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if err := s.UpdateUnitOfMeasure(tt.args.ctx, tt.args.uomId, tt.args.uom); err != nil && tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_service_GetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProdRepo := prodRepoMock.NewMockProductRepository(ctrl)
	ctx := context.Background()
	errMsg := "ERROR"
	err := errors.New(errMsg)
	productIDUuid := uuid.MustParse(productID)
	otherProductIDUuid := uuid.MustParse(otherProductID)
	otherProductIDUuid2 := uuid.MustParse(otherProductID2)

	mockProdRepo.EXPECT().GetProductById(ctx, otherProductIDUuid2).Return(nil, nil).AnyTimes()
	mockProdRepo.EXPECT().GetProductById(ctx, otherProductIDUuid).Return(nil, err).AnyTimes()
	mockProdRepo.EXPECT().GetProductById(ctx, productIDUuid).Return(&prodEntity.Product{
		Name: "mouse",
	}, nil).AnyTimes()
	type fields struct {
		productRepository *prodRepoMock.MockProductRepository
	}
	type args struct {
		ctx       context.Context
		productId uuid.UUID
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "GetProductById_Error_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:       ctx,
				productId: otherProductIDUuid,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.Internal, "Error when getting product by id :"+errMsg),
		},
		{
			name: "GetProductById_ProductIdNotFound_ReturnTheError",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:       ctx,
				productId: otherProductIDUuid2,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.NotFound, "Product id not found"),
		},
		{
			name: "GetProductById_NoError_Success",
			fields: fields{
				productRepository: mockProdRepo,
			},
			args: args{
				ctx:       ctx,
				productId: productIDUuid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, tt.fields.productRepository, nil, nil)
			if p, err := s.GetProductById(tt.args.ctx, tt.args.productId); err != nil && tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, p)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, p)
			}
		})
	}
}
