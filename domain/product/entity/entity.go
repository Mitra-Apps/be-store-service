package entity

import (
	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type ProductCategory struct {
	base_model.BaseModel
	Name         string         `gorm:"type:varchar(255);not null;unique"`
	IsActive     bool           `gorm:"type:bool;not null;default:TRUE"`
	ProductTypes []*ProductType `gorm:"foreignKey:ProductCategoryID"`
}

type ProductType struct {
	base_model.BaseModel
	Name              string     `gorm:"type:varchar(255);not null;unique"`
	IsActive          bool       `gorm:"type:bool;not null;default:TRUE"`
	ProductCategoryID uuid.UUID  `gorm:"type:uuid;not null"`
	Products          []*Product `gorm:"foreignKey:ProductTypeID"`
}

type UnitOfMeasure struct {
	base_model.BaseModel
	Name     string     `gorm:"type:varchar(255);not null;unique"`
	Symbol   string     `gorm:"type:varchar(50);not null;unique"`
	IsActive bool       `gorm:"type:bool;not null;default:TRUE"`
	Products []*Product `gorm:"foreignKey:UomID"`
}

type Product struct {
	base_model.BaseModel
	StoreID       uuid.UUID `gorm:"type:uuid;not null"`
	Name          string    `gorm:"type:varchar(255);not null;unique"`
	SaleStatus    bool      `gorm:"type:bool;not null;default:FALSE"`
	Price         float64   `gorm:"decimal(17,2); not null; default:0"`
	Stock         int64     `gorm:"type:int;"`
	UomID         uuid.UUID `gorm:"type:uuid;not null"`
	ProductTypeID uuid.UUID `gorm:"type:uuid;not null"`
}

func (p *Product) FromProto(product *pb.Product) error {
	if product.Id != "" {
		id, err := uuid.Parse(product.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product id")
		}
		p.ID = id
	}

	if product.StoreId != "" {
		storeId, err := uuid.Parse(product.StoreId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for store id")
		}
		p.StoreID = storeId
	}

	if product.UomId != "" {
		uomId, err := uuid.Parse(product.UomId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for unit of measure id")
		}
		p.UomID = uomId
	}

	if product.ProductTypeId != "" {
		prodTypeId, err := uuid.Parse(product.ProductTypeId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product type id")
		}
		p.ProductTypeID = prodTypeId
	}

	p.Name = product.Name
	p.SaleStatus = product.SaleStatus
	p.Price = product.Price
	p.Stock = product.Stock

	return nil
}
