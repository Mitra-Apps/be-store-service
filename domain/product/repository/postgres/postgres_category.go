package postgres

import (
	"context"
	"fmt"

	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/Mitra-Apps/be-store-service/types"
)

func (p *Postgres) GetCategoriesByStoreIdRepo(ctx context.Context, input types.GetCategoriesByStoreIdRepoInput) (output types.GetCategoriesByStoreIdRepoOutput, err error) {
	categories := []*entity.Category{}

	productTable := "products"
	productCategoryRelationsTable := "product_category_relations"
	aliasCategoryTable := "\"categories\""
	aliasProductTable := "\"products\""
	aliasProductCategoryRelationsTable := "\"product_category_relations\""
	selectCategoryColumn := fmt.Sprintf("%s.id, %s.name", aliasCategoryTable, aliasCategoryTable)

	tx := p.db.WithContext(ctx).
		Select(selectCategoryColumn).
		Joins(fmt.Sprintf(" INNER JOIN %s on %s.category_id = %s.id ", productCategoryRelationsTable, aliasProductCategoryRelationsTable, aliasCategoryTable)).
		Joins(fmt.Sprintf(" INNER JOIN %s on %s.product_id = %s.id ", productTable, aliasProductCategoryRelationsTable, aliasProductTable)).
		Where(fmt.Sprintf(" %s.store_id = ? ", aliasProductTable), input.StoreID)

	if !input.IsIncludeDeactivated {
		tx = tx.Where("is_active = ?", true)
		selectCategoryColumn = fmt.Sprintf("%s, %s.is_active", selectCategoryColumn, aliasCategoryTable)
	}

	tx = tx.Group(selectCategoryColumn).
		Order(fmt.Sprintf("%s.name ASC", aliasCategoryTable))

	if err := tx.Find(&categories).Error; err != nil {
		return output, err
	}

	output.Categories = categories

	return output, err
}
