package types

import (
	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	"github.com/google/uuid"
)

type GetProductsByStoreIdParams struct {
	Page                 int32
	Limit                int32
	StoreID              uuid.UUID
	ProductTypeId        *int64
	IsIncludeDeactivated bool
	OrderBy              string
	Direction            string
}

type GetProductsByStoreIdRepoParams struct {
	Pagination           base_model.Pagination
	StoreID              uuid.UUID
	ProductTypeId        *int64
	IsIncludeDeactivated bool
	OrderBy              string
	Direction            string
}
