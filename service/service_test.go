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
	userID       = "8b15140c-f6d0-4f2f-8302-57383a51adaf"
	otherUserID  = "2f27d467-9f83-4170-96ab-36e0994f37ca"
	storeID      = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
	otherStoreID = "52d11042-8c45-453e-86af-fe1e4d7facf6"
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
			s := New(tt.fields.storeRepository, nil, nil)
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
			service := New(storeRepository, nil, storage)

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
	products = append(products, &prodEntity.Product{
		StoreID: storeIdUuid,
		Name:    "indomie",
	})
	products = append(products, &prodEntity.Product{
		StoreID: storeIdUuid,
		Name:    "beras",
	})
	existedProducts := []*prodEntity.Product{}
	existedProducts = append(existedProducts, &prodEntity.Product{
		StoreID: storeIdUuid,
		Name:    "bakso aci",
	})
	existedProducts = append(existedProducts, &prodEntity.Product{
		StoreID: storeIdUuid,
		Name:    "keju",
	})
	roleNames := []string{"merchant"}
	adminRoleNames := []string{"merchant", "admin"}
	existedProdNames := []string{"bakso aci", "keju"}
	mockProdRepo.EXPECT().GetProductsByStoreIdAndNames(gomock.Any(), storeIdUuid, existedProdNames).Return(existedProducts, nil).AnyTimes()
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
			s := New(tt.fields.storeRepository, tt.fields.productRepository, nil)
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
			service := New(storeRepository, nil, storage)

			tc.setupMocks(storeRepository, storage)
			result, err := service.UpdateStore(ctx, tc.inputStore.storeID, tc.inputStore.store)
			assert.Equal(t, tc.expectedError, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
