package postgres

import (
	"context"
	"errors"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"
	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
			return nil, status.Errorf(codes.NotFound, "Store with id %s not found", storeID)
		}
		return nil, status.Errorf(codes.Internal, "Error when getting store :"+err.Error())
	}
	return &store, nil
}

func (p *postgres) UpdateStore(ctx context.Context, update *entity.Store) (*entity.Store, error) {
	if update.ID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "store id is required")
	}

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.Store{}).Where("id = ?", update.ID.String()).Updates(map[string]interface{}{
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
			if err == gorm.ErrRecordNotFound {
				return status.Errorf(codes.NotFound, "store with ID %s not found", update.ID.String())
			}

			return status.Errorf(codes.Internal, "failed to update store: %v", err)
		}

		if err := p.updateStoreTags(ctx, tx, update.ID, update.Tags); err != nil {
			return err
		}

		if err := p.updateStoreHours(ctx, tx, update.ID, update.Hours); err != nil {
			return err
		}

		if err := p.updateStoreImages(ctx, tx, update.ID, update.Images); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return p.GetStore(ctx, update.ID.String())
}

// updateStoreTags is a helper function to update the store tags.
func (p *postgres) updateStoreTags(ctx context.Context, tx *gorm.DB, storeID uuid.UUID, tags []*entity.StoreTag) error {
	exist, err := p.GetStore(ctx, storeID.String())
	if err != nil {
		return err
	}

	if err := tx.Model(exist).Association("Tags").Unscoped().Clear(); err != nil {
		return err
	}

	for _, tag := range tags {
		if err := tx.First(tag, "tag_name = ?", tag.TagName).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if tag.ID == uuid.Nil {
			if err := tx.Create(tag).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(exist).Association("Tags").Append(tag); err != nil {
			return err
		}
	}

	return nil
}

// updateStoreHours is a helper function to update the store hours.
func (p *postgres) updateStoreHours(ctx context.Context, tx *gorm.DB, storeID uuid.UUID, hours []*entity.StoreHour) error {
	if err := tx.Unscoped().Delete(hours, "store_id = ?", storeID.String()).Error; err != nil {
		return err
	}

	if len(hours) > 0 {
		return tx.Create(hours).Error
	}

	return nil
}

// updateStoreImages is a helper function to update the store images.
func (p *postgres) updateStoreImages(ctx context.Context, tx *gorm.DB, storeID uuid.UUID, images []*entity.StoreImage) error {
	if err := tx.Unscoped().Delete(images, "store_id = ?", storeID.String()).Error; err != nil {
		return err
	}

	if len(images) > 0 {
		return tx.Create(images).Error
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
