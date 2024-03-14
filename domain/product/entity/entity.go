package entity

import (
	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type ProductCategory struct {
	base_model.BaseMasterDataModel
	Name         string         `gorm:"type:varchar(255);not null;unique"`
	IsActive     bool           `gorm:"type:bool;not null"`
	ProductTypes []*ProductType `gorm:"foreignKey:ProductCategoryID"`
}

type ProductType struct {
	base_model.BaseMasterDataModel
	Name              string          `gorm:"type:varchar(255);not null"`
	IsActive          bool            `gorm:"type:bool;not null"`
	ProductCategoryID int64           `gorm:"type:uuid;not null"`
	ProductCategory   ProductCategory `gorm:"foreignKey:ProductCategoryID"`
}

type UnitOfMeasure struct {
	base_model.BaseMasterDataModel
	Name     string `gorm:"type:varchar(255);not null;unique"`
	Symbol   string `gorm:"type:varchar(50);not null;unique"`
	IsActive bool   `gorm:"type:bool;not null"`
}

type Product struct {
	base_model.BaseModel
	StoreID             uuid.UUID       `gorm:"type:uuid;not null"`
	Name                string          `gorm:"type:varchar(255);not null"`
	SaleStatus          bool            `gorm:"type:bool;not null"`
	Price               float64         `gorm:"decimal(17,2); not null; default:0"`
	Stock               int64           `gorm:"type:int;"`
	UomID               int64           `gorm:"type:int;not null"`
	ProductTypeID       int64           `gorm:"type:uuid;not null"`
	Images              []*ProductImage `gorm:"foreignKey:ProductId"`
	ProductType         ProductType     `gorm:"foreignKey:ProductTypeID"`
	ProductTypeName     string          `gorm:"-"`
	ProductCategoryID   int64           `gorm:"-"`
	ProductCategoryName string          `gorm:"-"`
}

type ProductImage struct {
	base_model.BaseModel
	ProductId      uuid.UUID `gorm:"type:uuid;not null"`
	ImageId        uuid.UUID `gorm:"type:uuid;not null"`
	ImageBase64Str string    `gorm:"-"`
	ImageURL       string    `gorm:"-"`
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

	for _, i := range product.Images {
		p.Images = append(p.Images, &ProductImage{
			ImageBase64Str: i.ImageBase64Str,
		})
	}

	p.Name = product.Name
	p.SaleStatus = product.SaleStatus
	p.Price = product.Price
	p.Stock = product.Stock
	p.UomID = product.UomId
	p.ProductTypeID = product.ProductTypeId
	p.ProductTypeName = product.ProductTypeName
	p.ProductCategoryID = product.ProductCategoryId
	p.ProductCategoryName = product.ProductCategoryName

	return nil
}

func (p *ProductCategory) FromProto(category *pb.ProductCategory) {
	p.ID = category.Id
	p.Name = category.Name
	p.IsActive = category.IsActive
}

func (p *ProductType) FromProto(prodType *pb.ProductType) error {
	p.ID = prodType.Id
	p.ProductCategoryID = prodType.ProductCategoryId
	p.Name = prodType.Name
	p.IsActive = prodType.IsActive

	return nil
}

func (u *UnitOfMeasure) FromProto(uom *pb.UnitOfMeasure) error {
	u.ID = uom.Id
	u.Name = uom.Name
	u.Symbol = uom.Symbol
	u.IsActive = uom.IsActive

	return nil
}

func (p *Product) ToProto() *pb.Product {
	if p == nil {
		return nil
	}
	images := []*pb.ProductImage{}
	for _, i := range p.Images {
		images = append(images, &pb.ProductImage{
			Id:       i.ID.String(),
			ImageId:  i.ImageId.String(),
			ImageUrl: i.ImageURL,
		})
	}
	return &pb.Product{
		Id:                  p.ID.String(),
		StoreId:             p.StoreID.String(),
		Name:                p.Name,
		SaleStatus:          p.SaleStatus,
		Price:               p.Price,
		Stock:               p.Stock,
		UomId:               p.UomID,
		ProductTypeId:       p.ProductTypeID,
		ProductTypeName:     p.ProductTypeName,
		ProductCategoryId:   p.ProductCategoryID,
		ProductCategoryName: p.ProductCategoryName,
		Images:              images,
	}
}

func (u *UnitOfMeasure) ToProto() *pb.UnitOfMeasure {
	if u == nil {
		return nil
	}
	return &pb.UnitOfMeasure{
		Id:       u.ID,
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
		Id:       c.ID,
		Name:     c.Name,
		IsActive: c.IsActive,
	}
}

func (t *ProductType) ToProto() *pb.ProductType {
	if t == nil {
		return nil
	}
	return &pb.ProductType{
		Id:                t.ID,
		Name:              t.Name,
		IsActive:          t.IsActive,
		ProductCategoryId: t.ProductCategoryID,
	}
}
