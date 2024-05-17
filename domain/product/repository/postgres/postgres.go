package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/Mitra-Apps/be-store-service/types"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Postgres struct {
	db  *gorm.DB
	trx *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{db, nil}
}

func (p *Postgres) GetProductsByStoreId(ctx context.Context, params types.GetProductsByStoreIdRepoParams) ([]*entity.Product, base_model.Pagination, error) {
	prods := []*entity.Product{}

	paging := base_model.Pagination{
		Page:  params.Pagination.Page,
		Limit: params.Pagination.Limit,
	}

	tx := p.db.WithContext(ctx).
		Preload("Images").
		Preload("ProductType").
		Preload("ProductType.ProductCategory").
		Where("store_id = ?", params.StoreID)
	if !params.IsIncludeDeactivated {
		tx = tx.Where("sale_status = ?", true)
	}

	if params.ProductTypeId != nil {
		tx = tx.Where("product_type_id = ?", *params.ProductTypeId)
	}

	tx = tx.Order(fmt.Sprintf("%s %s", params.OrderBy, params.Direction))
	err := tx.
		Scopes(p.paginate(prods, &paging, tx, int64(len(prods)))).
		Find(&prods).
		Error

	if err != nil {
		if strings.Contains(err.Error(), ErrIncorrectSqlSyntax) || strings.Contains(err.Error(), ErrInvalidColumnName) {
			return nil, paging, fmt.Errorf("incorrect orderBy or direction value")
		}
		if errors.Is(err, gorm.ErrInvalidField) {
			return nil, paging, nil
		}
		return nil, paging, err
	}

	for _, p := range prods {
		p.ProductTypeName = p.ProductType.Name
		p.ProductCategoryID = p.ProductType.ProductCategoryID
		p.ProductCategoryName = p.ProductType.ProductCategory.Name
	}

	return prods, paging, nil
}

func (p *Postgres) GetProductById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var prod entity.Product
	tx := p.db.WithContext(ctx).
		Preload("Images").
		Preload("ProductType").
		Preload("ProductType.ProductCategory").
		First(&prod, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	prod.ProductTypeName = prod.ProductType.Name
	prod.ProductCategoryID = prod.ProductType.ProductCategoryID
	prod.ProductCategoryName = prod.ProductType.ProductCategory.Name

	return &prod, nil
}

func (p *Postgres) GetProductsByStoreIdAndNames(ctx context.Context, storeID uuid.UUID, names []string) ([]*entity.Product, error) {
	lowerCaseNames := []string{}
	for _, s := range names {
		lowerCaseNames = append(lowerCaseNames, strings.ToLower(s))
	}
	prods := []*entity.Product{}
	tx := p.db.WithContext(ctx).
		Preload("Images").
		Where("store_id = ? AND LOWER(name) IN ?", storeID, lowerCaseNames).Find(&prods)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return prods, nil
}

func (p *Postgres) InitiateTransaction(ctx context.Context) bool {
	newTrx := p.trx == nil
	if newTrx {
		p.trx = p.db.WithContext(ctx).Begin()
	}
	return newTrx
}

func (p *Postgres) TransactionRollback() {
	p.trx.Rollback()
	p.trx = nil
}

func (p *Postgres) TransactionCommit() error {
	defer func() {
		p.trx = nil
	}()

	if p.trx == nil {
		return status.Errorf(codes.Internal, "No transaction committed")
	}

	if err := p.trx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpsertProducts(ctx context.Context, products []*entity.Product) error {
	if err := p.trx.Omit("Images").
		Save(products).Error; err != nil {
		p.trx.Rollback()
		return err
	}

	return nil
}

func (p *Postgres) UpsertProductImages(ctx context.Context, productImages []*entity.ProductImage) error {
	funcScopeTrx := p.InitiateTransaction(ctx)

	if err := p.trx.Save(productImages).Error; err != nil {
		p.trx.Rollback()
		return err
	}

	if funcScopeTrx {
		return p.TransactionCommit()
	}

	return nil
}

func (p *Postgres) DeleteProductImages(ctx context.Context, productImages []*entity.ProductImage) error {
	funcScopeTrx := p.InitiateTransaction(ctx)

	if err := p.trx.Delete(productImages).Error; err != nil {
		p.trx.Rollback()
		return err
	}

	if funcScopeTrx {
		return p.TransactionCommit()
	}

	return nil
}

func (p *Postgres) DeleteProductById(ctx context.Context, id uuid.UUID) error {
	if err := p.trx.Where("id = ?", id.String()).Delete(&entity.Product{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetProductImagesByProductIds(ctx context.Context, productIds []uuid.UUID) ([]*entity.ProductImage, map[uuid.UUID][]*entity.ProductImage, error) {
	prodImages := []*entity.ProductImage{}
	if tx := p.db.WithContext(ctx).Where("product_id IN ?", productIds).Find(&prodImages); tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, tx.Error
	}

	byProduct := make(map[uuid.UUID][]*entity.ProductImage)
	for _, i := range prodImages {
		byProduct[i.ProductId] = append(byProduct[i.ProductId], i)
	}

	return prodImages, byProduct, nil
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

func (p *Postgres) GetUnitOfMeasureById(ctx context.Context, uomId int64) (*entity.UnitOfMeasure, error) {
	uom := entity.UnitOfMeasure{}
	err := p.db.WithContext(ctx).First(&uom, uomId).Error
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

func (p *Postgres) GetProductCategories(ctx context.Context, includeDeactivated bool) ([]*entity.ProductCategory, error) {
	categories := []*entity.ProductCategory{}
	tx := p.db.WithContext(ctx).Where("is_active = ?", includeDeactivated).Order("name ASC")
	if err := tx.Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return categories, nil
}

func (p *Postgres) GetProductCategoryByName(ctx context.Context, name string) (*entity.ProductCategory, error) {
	var category *entity.ProductCategory
	if err := p.db.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(name)).First(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (p *Postgres) GetProductCategoryById(ctx context.Context, id int64) (*entity.ProductCategory, error) {
	var category *entity.ProductCategory
	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (p *Postgres) UpsertProductCategory(ctx context.Context, prodCategory *entity.ProductCategory) error {
	return p.db.WithContext(ctx).Save(prodCategory).Error
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
	err := p.db.WithContext(ctx).
		Where("id IN ?", typeIds).
		Find(&prodTypes).Error
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

func (p *Postgres) paginate(value interface{}, pagination *base_model.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = int32(totalRecords)
	pagination.TotalPage = int32(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int32(pagination.Limit*(pagination.Page-1)) + int32(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(pagination.GetOffset())).Limit(int(pagination.GetLimit()))
	}
}

func (p *Postgres) GetProductTypeById(ctx context.Context, id int64) (*entity.ProductType, error) {
	var productType *entity.ProductType
	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&productType).Error; err != nil {
		return nil, err
	}
	return productType, nil
}
