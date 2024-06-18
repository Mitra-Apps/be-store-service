package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/Mitra-Apps/be-store-service/types"
	"gorm.io/gorm"
)

func (p *Postgres) GetCategoriesByStoreIdRepo(ctx context.Context, input types.GetCategoriesByStoreIdRepoInput) (output types.GetCategoriesByStoreIdRepoOutput, err error) {
	categories := []*entity.Category{}

	productTable := "products"
	productCategoryRelationsTable := "product_category_relations"
	aliasCategoryTable := "\"categories\""
	aliasProductTable := "\"products\""
	aliasProductCategoryRelationsTable := "\"product_category_relations\""
	selectCategoryColumn := fmt.Sprintf("%s.id, %s.name, %s.is_active", aliasCategoryTable, aliasCategoryTable, aliasCategoryTable)

	tx := p.db.WithContext(ctx).
		Select(selectCategoryColumn).
		Joins(fmt.Sprintf(" INNER JOIN %s on %s.category_id = %s.id ", productCategoryRelationsTable, aliasProductCategoryRelationsTable, aliasCategoryTable)).
		Joins(fmt.Sprintf(" INNER JOIN %s on %s.product_id = %s.id ", productTable, aliasProductCategoryRelationsTable, aliasProductTable)).
		Where(fmt.Sprintf(" %s.store_id = ? ", aliasProductTable), input.StoreID).
		Where(fmt.Sprintf("%s.is_active = ?", aliasCategoryTable), !input.IsIncludeDeactivated).
		Group(selectCategoryColumn).
		Order(fmt.Sprintf("%s.name ASC", aliasCategoryTable))

	if err := tx.Find(&categories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return output, nil
		}
		return output, err
	}

	output.Categories = categories

	return output, err
}
