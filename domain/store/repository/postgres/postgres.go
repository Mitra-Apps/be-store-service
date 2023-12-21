package postgres

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"

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
	if err := p.db.WithContext(ctx).Preload("Hours").Preload("Images").Preload("Tags").Where("id = ?", storeID).First(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (p *postgres) UpdateStore(ctx context.Context, storeID string, update *entity.Store) (*entity.Store, error) {
	tx := p.db.WithContext(ctx).Begin()

	var existingStore entity.Store
	if err := tx.Where("id = ?", storeID).Preload("Hours").Preload("Images").Preload("Tags").First(&existingStore).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&existingStore).Updates(map[string]interface{}{
		"store_name":   update.StoreName,
		"address":      update.Address,
		"city":         update.City,
		"state":        update.State,
		"zip_code":     update.ZipCode,
		"phone":        update.Phone,
		"email":        update.Email,
		"website":      update.Website,
		"map_location": update.MapLocation,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := p.updateStoreHours(ctx, existingStore, update.Hours); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := p.updateStoreImages(ctx, existingStore, update.Images); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := p.updateStoreTags(ctx, existingStore, update.Tags); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return update, nil
}

func (p *postgres) updateStoreHours(ctx context.Context, existingStore entity.Store, updatedHours []entity.StoreHour) error {
	for _, updatedHour := range updatedHours {
		var existingHour entity.StoreHour
		if err := p.db.WithContext(ctx).Where("id = ?", updatedHour.ID).First(&existingHour).Error; err != nil {
			return err
		}

		if err := p.db.WithContext(ctx).Model(&existingHour).Updates(map[string]interface{}{
			"day_of_week": updatedHour.DayOfWeek,
			"open":        updatedHour.Open,
			"close":       updatedHour.Close,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *postgres) updateStoreImages(ctx context.Context, existingStore entity.Store, updatedImages []entity.StoreImage) error {
	for _, existingImage := range existingStore.Images {
		found := false
		for _, updatedImage := range updatedImages {
			if existingImage.ID == updatedImage.ID {
				found = true
				break
			}
		}
		if !found {
			if err := p.db.WithContext(ctx).Delete(&existingImage).Error; err != nil {
				return err
			}
		}
	}

	for _, updatedImage := range updatedImages {
		var existingImage entity.StoreImage
		if updatedImage.ID.String() != "" {
			if err := p.db.WithContext(ctx).Where("id = ?", updatedImage.ID).First(&existingImage).Error; err != nil {
				return err
			}
		}

		if err := p.db.WithContext(ctx).Model(&existingImage).Updates(map[string]interface{}{
			"image_type": updatedImage.ImageType,
			"image_url":  updatedImage.ImageURL,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *postgres) updateStoreTags(ctx context.Context, existingStore entity.Store, updatedTags []entity.StoreTag) error {
	for _, existingCategory := range existingStore.Tags {
		found := false
		for _, updatedCategory := range updatedTags {
			if existingCategory.ID == updatedCategory.ID {
				found = true
				break
			}
		}
		if !found {
			if err := p.db.WithContext(ctx).Delete(&existingCategory).Error; err != nil {
				return err
			}
		}
	}

	for _, updatedCategory := range updatedTags {
		var existingCategory entity.StoreTag
		if updatedCategory.ID.String() != "" {
			if err := p.db.WithContext(ctx).Where("id = ?", updatedCategory.ID).First(&existingCategory).Error; err != nil {
				return err
			}
		}

		if err := p.db.WithContext(ctx).Model(&existingCategory).Updates(map[string]interface{}{
			"category_name": updatedCategory.TagName,
		}).Error; err != nil {
			return err
		}
	}
	return nil
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
