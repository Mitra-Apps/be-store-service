package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{db}
}

func (p *Postgres) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID, productTypeId *int64, isIncludeDeactivated bool) ([]*entity.Product, error) {
	prods := []*entity.Product{}
	tx := p.db.WithContext(ctx).Where("store_id = ?", storeID)
	if !isIncludeDeactivated {
		tx = tx.Where("sale_status = ?", true)
	}
	if productTypeId != nil {
		tx = tx.Where("product_type_id = ?", *productTypeId)
	}
	tx = tx.Order("name ASC")
	err := tx.Find(&prods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return prods, nil
}

func (p *Postgres) GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var prod entity.Product
	tx := p.db.WithContext(ctx).First(&prod, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &prod, nil
}

func (p *Postgres) GetProductsByStoreIdAndNames(ctx context.Context, storeID uuid.UUID, names []string) ([]*entity.Product, error) {
	lowerCaseNames := []string{}
	for _, s := range names {
		lowerCaseNames = append(lowerCaseNames, strings.ToLower(s))
	}
	prods := []*entity.Product{}
	tx := p.db.WithContext(ctx).Where("store_id = ? AND LOWER(name) IN ?", storeID, lowerCaseNames).Find(&prods)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return prods, nil
}

func (p *Postgres) UpsertProducts(ctx context.Context, products []*entity.Product) error {
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

func (p *Postgres) GetUnitOfMeasureByName(ctx context.Context, name string) (*entity.UnitOfMeasure, error) {
	uom := entity.UnitOfMeasure{}
	err := p.db.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(name)).First(&uom).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &uom, nil
}

func (p *Postgres) GetUnitOfMeasureBySymbol(ctx context.Context, symbol string) (*entity.UnitOfMeasure, error) {
	uom := entity.UnitOfMeasure{}
	err := p.db.WithContext(ctx).Where("symbol = ?", symbol).First(&uom).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &uom, nil
}

func (p *Postgres) GetUnitOfMeasures(ctx context.Context, isIncludeDeactivated bool) ([]*entity.UnitOfMeasure, error) {
	uom := []*entity.UnitOfMeasure{}
	var err error
	tx := p.db.WithContext(ctx)
	if !isIncludeDeactivated {
		tx = tx.Where("is_active = ?", true)
	}
	tx = tx.Order("name ASC")
	err = tx.Find(&uom).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return uom, nil
}

func (p *Postgres) GetUnitOfMeasuresByIds(ctx context.Context, uomIds []int64) ([]*entity.UnitOfMeasure, error) {
	uoms := []*entity.UnitOfMeasure{}
	err := p.db.WithContext(ctx).Where("id IN ?", uomIds).Find(&uoms).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return uoms, nil
}

func (p *Postgres) UpsertUnitOfMeasure(ctx context.Context, uom *entity.UnitOfMeasure) error {
	tx := p.db.WithContext(ctx).Begin()

	if err := tx.Save(uom).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetProductCategories(ctx context.Context, isIncludeDeactivated bool) ([]*entity.ProductCategory, error) {
	cat := []*entity.ProductCategory{}
	var err error
	tx := p.db.WithContext(ctx)
	if !isIncludeDeactivated {
		tx = tx.Where("is_active = ?", true)
	}
	tx = tx.Order("name ASC")
	err = tx.Find(&cat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return cat, nil
}

func (p *Postgres) GetProductCategoryByName(ctx context.Context, name string) (*entity.ProductCategory, error) {
	cat := entity.ProductCategory{}
	err := p.db.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(name)).First(&cat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

func (p *Postgres) GetProductCategoryById(ctx context.Context, id int64) (*entity.ProductCategory, error) {
	cat := entity.ProductCategory{}
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

func (p *Postgres) UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error {
	tx := p.db.WithContext(ctx).Begin()

	if err := tx.Save(prodCategory).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetProductTypeByName(ctx context.Context, productCategoryID int64, name string) (*entity.ProductType, error) {
	prodType := entity.ProductType{}
	err := p.db.WithContext(ctx).Where("product_category_id = ? AND LOWER(name) = ?", productCategoryID, strings.ToLower(name)).First(&prodType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prodType, nil
}

func (p *Postgres) GetProductTypes(ctx context.Context, productCategoryID int64, isIncludeDeactivated bool) ([]*entity.ProductType, error) {
	types := []*entity.ProductType{}
	var err error
	tx := p.db.WithContext(ctx).Where("product_category_id = ?", productCategoryID)
	if !isIncludeDeactivated {
		tx = tx.Where("is_active = ?", true)
	}
	tx = tx.Order("name ASC")
	err = tx.Find(&types).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return types, nil
}

func (p *Postgres) GetProductTypesByIds(ctx context.Context, typeIds []int64) ([]*entity.ProductType, error) {
	prodTypes := []*entity.ProductType{}
	err := p.db.WithContext(ctx).Where("id IN ?", typeIds).Find(&prodTypes).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return prodTypes, nil
}

func (p *Postgres) UpsertProductType(ctx context.Context, prodType *entity.ProductType) error {
	tx := p.db.WithContext(ctx).Begin()

	if err := tx.Save(prodType).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
