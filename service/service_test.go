package service

import (
	"context"
	"testing"

	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	storeRepoMock "github.com/Mitra-Apps/be-store-service/domain/store/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	storeID = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
)

func Test_service_OpenCloseStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStoreRepo := storeRepoMock.NewMockStoreServiceRepository(ctrl)
	ctx := context.Background()
	storeIdUuid, _ := uuid.Parse(storeID)
	mockStoreRepo.EXPECT().OpenCloseStore(ctx, storeIdUuid, false).Return(nil)

	type fields struct {
		storeRepository   *storeRepoMock.MockStoreServiceRepository
		storageRepository repository.Storage
	}
	type args struct {
		ctx      context.Context
		storeID  string
		isActive bool
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
				ctx:      ctx,
				storeID:  "",
				isActive: false,
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
				ctx:      ctx,
				storeID:  "aaa",
				isActive: false,
			},
			wantErr:       true,
			expectedError: status.Errorf(codes.InvalidArgument, "store id should be uuid"),
		},
		{
			name: "OpenCloseStore_NoError_ReturnNil",
			fields: fields{
				storeRepository: mockStoreRepo,
			},
			args: args{
				ctx:      ctx,
				storeID:  storeID,
				isActive: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.storeRepository, nil)
			if err := s.OpenCloseStore(tt.args.ctx, tt.args.storeID, tt.args.isActive); tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
