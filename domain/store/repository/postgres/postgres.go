package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) repository.StoreServiceRepository {
	return &postgres{db}
}

func (p *postgres) CreateStore(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	tx := p.db.WithContext(ctx).Begin()

	if err := tx.Save(store).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return store, nil
}

func (p *postgres) GetStore(ctx context.Context, storeID string) (*entity.Store, error) {
	var store entity.Store
	err := p.db.WithContext(ctx).
		Model(&store).
		Preload("Hours").
		Preload("Images").
		Preload("Tags").
		Where("id = ?", storeID).
		First(&store).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("store with ID %s not found", storeID)
		}
		return nil, fmt.Errorf("error retrieving store: %w", err)
	}
	return &store, nil
}

func (p *postgres) UpdateStore(ctx context.Context, storeID string, update *entity.Store) error {
	exist, err := p.GetStore(ctx, storeID)
	if err != nil {
		return err
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Model(exist).Updates(map[string]interface{}{
			"user_id":           update.UserID,
			"store_name":        update.StoreName,
			"store_description": update.StoreDescription,
			"address":           update.Address,
			"city":              update.City,
			"state":             update.State,
			"zip_code":          update.ZipCode,
			"phone":             update.Phone,
			"email":             update.Email,
			"website":           update.Website,
			"location_lat":      update.LocationLat,
			"location_lng":      update.LocationLng,
			"status":            update.Status,
			"is_active":         update.IsActive,
		}).Error
		if err != nil {
			return err
		}

		exist.Tags = nil
		exist.Tags = append(exist.Tags, update.Tags...)

		exist.Hours = nil
		exist.Hours = append(exist.Hours, update.Hours...)

		exist.Images = nil
		exist.Images = append(exist.Images, update.Images...)

		return tx.Save(exist).Error
	})
}

func (p *postgres) DeleteStore(ctx context.Context, storeID string) error {
	if err := p.db.WithContext(ctx).Where("id = ?", storeID).Delete(&entity.Store{}).Error; err != nil {
		return err
	}
	return nil
}

func (p *postgres) ListStores(ctx context.Context) ([]*entity.Store, error) {
	var stores []*entity.Store
	if err := p.db.WithContext(ctx).Preload("Hours").Preload("Images").Preload("Tags").Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (p *postgres) OpenCloseStore(ctx context.Context, storeId uuid.UUID, isActive bool) error {
	tx := p.db.WithContext(ctx).Model(entity.Store{}).Where("id = ?", storeId).
		UpdateColumn("is_active", isActive)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// GetStoreByUserID retrieves a store by its user ID.
func (p *postgres) GetStoreByUserID(ctx context.Context, userID uuid.UUID) (*entity.Store, error) {
	var store entity.Store
	if err := p.db.WithContext(ctx).
		Preload("Hours").
		Preload("Images").
		Preload("Tags").
		First(&store, "stores.user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &store, nil
}
