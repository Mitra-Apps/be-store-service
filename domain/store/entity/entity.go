package entity

import (
	"time"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/google/uuid"
)

type Store struct {
	Id          uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserId      uuid.UUID  `gorm:"type:uuid;not null"`
	Name        string     `gorm:"type:varchar(255);not null"`
	Address     string     `gorm:"type:text;null"`
	MapLocation string     `gorm:"type:varchar(255);null"`
	LogoImageId *uuid.UUID `gorm:"type:varchar(255);null"`
	Status      string     `gorm:"type:varchar(50);not null"`
	IsActive    bool       `gorm:"type:bool;not null"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz;null"`
	UpdatedBy   *uuid.UUID `gorm:"type:uuid;null"`
}

func (s *Store) ToProto() *pb.Store {
	var logoImageId string
	if s.LogoImageId != nil {
		logoImageId = s.LogoImageId.String()
	}
	return &pb.Store{
		Id:          s.Id.String(),
		UserId:      s.UserId.String(),
		Name:        s.Name,
		Address:     s.Address,
		MapLocation: s.MapLocation,
		LogoImageId: logoImageId,
		Status:      s.Status,
		IsActive:    s.IsActive,
	}
}
