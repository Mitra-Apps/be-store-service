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
	Name              string     `gorm:"type:varchar(255);not null"`
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
	Name          string    `gorm:"type:varchar(255);not null"`
	SaleStatus    bool      `gorm:"type:bool;not null;default:TRUE"`
	Price         float64   `gorm:"decimal(17,2); not null; default:0"`
	Stock         int64     `gorm:"type:int;"`
	UomID         uuid.UUID `gorm:"type:uuid;not null"`
	ProductTypeID uuid.UUID `gorm:"type:uuid;not null"`
}

func (p *Product) FromProto(product *pb.Product, storeIdPrm *string) error {
	if product.Id != "" {
		id, err := uuid.Parse(product.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product id")
		}
		p.ID = id
	}

	storeId := product.StoreId
	if storeIdPrm != nil {
		storeId = *storeIdPrm
	}

	if product.StoreId != "" {
		storeId, err := uuid.Parse(storeId)
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

func (p *ProductCategory) FromProto(category *pb.ProductCategory) error {
	if category.Id != "" {
		id, err := uuid.Parse(category.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product category id")
		}
		p.ID = id
	}

	p.Name = category.Name
	p.IsActive = category.IsActive

	return nil
}

func (p *ProductType) FromProto(prodType *pb.ProductType) error {
	if prodType.Id != "" {
		id, err := uuid.Parse(prodType.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product type id")
		}
		p.ID = id
	}

	if prodType.ProductCategoryId != "" {
		categoryId, err := uuid.Parse(prodType.ProductCategoryId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product category id")
		}
		p.ProductCategoryID = categoryId
	}

	p.Name = prodType.Name
	p.IsActive = prodType.IsActive

	return nil
}

func (u *UnitOfMeasure) FromProto(uom *pb.UnitOfMeasure) error {
	if uom.Id != "" {
		uomId, err := uuid.Parse(uom.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid uuid for product category id")
		}
		u.ID = uomId
	}

	u.Name = uom.Name
	u.Symbol = uom.Symbol
	u.IsActive = uom.IsActive

	return nil
}

func (p *Product) ToProto() *pb.Product {
	if p == nil {
		return nil
	}
	return &pb.Product{
		Id:            p.ID.String(),
		StoreId:       p.StoreID.String(),
		Name:          p.Name,
		SaleStatus:    p.SaleStatus,
		Price:         p.Price,
		Stock:         p.Stock,
		UomId:         p.UomID.String(),
		ProductTypeId: p.ProductTypeID.String(),
	}
}

func (u *UnitOfMeasure) ToProto() *pb.UnitOfMeasure {
	if u == nil {
		return nil
	}
	return &pb.UnitOfMeasure{
		Id:       u.ID.String(),
		Name:     u.Name,
		Symbol:   u.Symbol,
		IsActive: u.IsActive,
	}
}

func (c *ProductCategory) ToProto() *pb.ProductCategory {
	if c == nil {
		return nil
	}
	return &pb.ProductCategory{
		Id:       c.ID.String(),
		Name:     c.Name,
		IsActive: c.IsActive,
	}
}

func (t *ProductType) ToProto() *pb.ProductType {
	if t == nil {
		return nil
	}
	return &pb.ProductType{
		Id:                t.ID.String(),
		Name:              t.Name,
		IsActive:          t.IsActive,
		ProductCategoryId: t.ProductCategoryID.String(),
	}
}
