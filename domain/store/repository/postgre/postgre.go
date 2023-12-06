package postgre

import (
	"context"

	"github.com/Mitra-Apps/be-store-service/domain/store/entity"

	"gorm.io/gorm"
)

type Postgre struct {
	db *gorm.DB
}

func NewPostgre(db *gorm.DB) *Postgre {
	return &Postgre{db}
}

func (p *Postgre) GetAll(ctx context.Context) ([]*entity.Store, error) {
	var stores []*entity.Store
	res := p.db.Order("created_at DESC").Find(&stores)
	if res.Error == gorm.ErrEmptySlice || res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return stores, nil
}
