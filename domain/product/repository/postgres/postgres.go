package postgres

import (
	"context"
	"errors"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) *postgres {
	return &postgres{db}
}

func (p *postgres) GetProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	prods := []*entity.Product{}
	tx := p.db.WithContext(ctx).Find(&prods, ids)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return prods, nil
}

func (p *postgres) CreateProducts(ctx context.Context, products []entity.Product) error {
	tx := p.db.WithContext(ctx).Begin()

	if err := tx.Save(products).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
