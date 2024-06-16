package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	storeID = "7d56be32-70a2-4f49-b66b-63e6f8e719d5"
)

var storeUUID, _ = uuid.Parse(storeID)

func TestCreateStore(t *testing.T) {
	if *integrationTest {
		storeRepo := NewPostgres(db)
		ctx := context.Background()
		store := &entity.Store{
			BaseModel: base_model.BaseModel{
				ID: storeUUID,
			},
			StoreName:        "Test Store",
			UserID:           uuid.New(),
			StoreDescription: "test desc",
			Address:          "test address",
		}

		t.Run("Success", func(t *testing.T) {
			result, err := storeRepo.CreateStore(ctx, store)
			assert.Nil(t, err)
			assert.Equal(t, store, result)
		})

		t.Run("Save Error", func(t *testing.T) {
			result, _ := storeRepo.CreateStore(ctx, store)
			result.BaseModel.ID = uuid.Nil
			result, err = storeRepo.CreateStore(ctx, result)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestGetStore(t *testing.T) {
	if *integrationTest {
		storeRepo := NewPostgres(db)
		ctx := context.Background()

		// Insert test data
		store := &entity.Store{
			BaseModel: base_model.BaseModel{
				ID: storeUUID,
			},
			StoreName:        "Test Store",
			UserID:           uuid.New(),
			StoreDescription: "test desc",
			Address:          "test address",
		}
		db.Create(store)

		t.Run("Success", func(t *testing.T) {
			result, err := storeRepo.GetStore(ctx, storeID)
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, err := storeRepo.GetStore(ctx, uuid.NewString())
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, status.Errorf(codes.NotFound, "Store is not found"), err)
		})

		t.Run("Other error", func(t *testing.T) {
			result, err := storeRepo.GetStore(ctx, "123")
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestUpdateStore(t *testing.T) {
	storeRepo := NewPostgres(db)
	ctx := context.Background()

	// Insert test data
	initialStore := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeUUID,
		},
		StoreName:        "Test Store",
		UserID:           uuid.New(),
		StoreDescription: "test desc",
		Address:          "test address",
	}
	db.Create(initialStore)

	// Update data
	updatedStore := &entity.Store{
		BaseModel: base_model.BaseModel{
			ID: storeUUID,
		},
		StoreName:        "Updated Store",
		UserID:           uuid.New(),
		StoreDescription: "Updated Description",
		Address:          "Updated Address",
		City:             "Updated City",
		State:            "Updated State",
		ZipCode:          "12345",
		Phone:            "123-456-7890",
		Email:            "updated@example.com",
		Website:          "http://updated.example.com",
		LocationLat:      40.7128,
		LocationLng:      -74.0060,
		Status:           "Active",
		IsActive:         true,
		Hours:            []*entity.StoreHour{{StoreID: storeUUID, DayOfWeek: "Monday", Open: "08:00", Close: "17:00"}},
		Images:           []*entity.StoreImage{{StoreID: storeUUID, ImageURL: "updated_image.jpg"}},
	}

	t.Run("Success", func(t *testing.T) {
		result, err := storeRepo.UpdateStore(ctx, updatedStore)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedStore.StoreName, result.StoreName)
		assert.Equal(t, updatedStore.StoreDescription, result.StoreDescription)
		assert.Equal(t, updatedStore.Address, result.Address)
		assert.Equal(t, updatedStore.City, result.City)
		assert.Equal(t, updatedStore.State, result.State)
		assert.Equal(t, updatedStore.ZipCode, result.ZipCode)
		assert.Equal(t, updatedStore.Phone, result.Phone)
		assert.Equal(t, updatedStore.Email, result.Email)
		assert.Equal(t, updatedStore.Website, result.Website)
		assert.Equal(t, updatedStore.LocationLat, result.LocationLat)
		assert.Equal(t, updatedStore.LocationLng, result.LocationLng)
		assert.Equal(t, updatedStore.Status, result.Status)
		assert.Equal(t, updatedStore.IsActive, result.IsActive)
		assert.Equal(t, updatedStore.Hours[0].Open, result.Hours[0].Open)
		assert.Equal(t, updatedStore.Images[0].ImageURL, result.Images[0].ImageURL)
	})

	t.Run("InvalidID", func(t *testing.T) {
		invalidStore := &entity.Store{
			BaseModel: base_model.BaseModel{
				ID: uuid.Nil,
			}}
		result, err := storeRepo.UpdateStore(ctx, invalidStore)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "store id is required")
	})

	t.Run("StoreNotFound", func(t *testing.T) {
		nonExistentStore := &entity.Store{
			BaseModel: base_model.BaseModel{
				ID: uuid.New(),
			},
			StoreName:        "Non-Existent Store",
			StoreDescription: "Non-Existent Description",
			Address:          "Non-Existent Address",
			UserID:           uuid.New(),
		}
		result, err := storeRepo.UpdateStore(ctx, nonExistentStore)
		fmt.Println("error is : ", err.Error())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}
