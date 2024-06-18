package types

import (
	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	"github.com/Mitra-Apps/be-store-service/domain/product/entity"
	"github.com/google/uuid"
)

type GetProductsByStoreIdParams struct {
	Page                 int32
	Limit                int32
	StoreID              uuid.UUID
	ProductTypeId        *int64
	ProductCategoryId    []int64
	IsIncludeDeactivated bool
	OrderBy              string
	Direction            string
	Search               *string
	UserID               uuid.UUID
}

type GetProductsByStoreIdRepoParams struct {
	Pagination           base_model.Pagination
	StoreID              uuid.UUID
	ProductTypeId        *int64
	ProductCategoryId    []int64
	IsIncludeDeactivated bool
	OrderBy              string
	Direction            string
	Search               *string
}

type GetProductCategoriesByStoreIdParams struct {
	IsIncludeDeactivated bool
	StoreID              uuid.UUID
	UserID               uuid.UUID
}

type GetCategoriesByStoreIdInput struct {
	IsIncludeDeactivated bool
	StoreID              uuid.UUID
	UserID               uuid.UUID
}

type GetCategoriesByStoreIdOutput struct {
	Categories []*entity.Category
}

type GetCategoriesByStoreIdRepoInput struct {
	IsIncludeDeactivated bool
	StoreID              uuid.UUID
}

type GetCategoriesByStoreIdRepoOutput struct {
	Categories []*entity.Category
}

type SaveCategoriesByStoreIdToRedisInput struct {
	IsIncludeDeactivated bool
	StoreID              uuid.UUID
}

type SaveCategoriesByStoreIdToRedisOutput struct {
	Categories []*entity.Category
}

type JwtCustomClaim struct {
	Roles            []string         `json:"roles"`
	RegisteredClaims RegisteredClaims `json:"registered_claims"`
}

type RegisteredClaims struct {
	Subject string `json:"sub"`
}
