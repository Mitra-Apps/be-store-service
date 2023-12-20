package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DayOfWeekEnum represents the days of the week.
type DayOfWeekEnum string

const (
	Monday    DayOfWeekEnum = "Monday"
	Tuesday   DayOfWeekEnum = "Tuesday"
	Wednesday DayOfWeekEnum = "Wednesday"
	Thursday  DayOfWeekEnum = "Thursday"
	Friday    DayOfWeekEnum = "Friday"
	Saturday  DayOfWeekEnum = "Saturday"
	Sunday    DayOfWeekEnum = "Sunday"
)

// BaseModel contains common fields for all models.
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedBy uuid.UUID
}

// Store represents a store model.
type Store struct {
	BaseModel
	StoreName   string `gorm:"not null;unique"`
	Address     string `gorm:"not null,type:text"`
	City        string `gorm:"not null"`
	State       string `gorm:"not null"`
	ZipCode     string `gorm:"not null"`
	Phone       string `gorm:"not null"`
	Email       string `gorm:"not null;unique"`
	Website     string
	MapLocation string
	Tags        []StoreTag   `gorm:"many2many:store_store_tags"`
	Hours       []StoreHour  `gorm:"foreignKey:StoreID"`
	Images      []StoreImage `gorm:"foreignKey:StoreID"`
}

// StoreImage represents an image associated with a store.
type StoreImage struct {
	BaseModel
	StoreID   uuid.UUID `gorm:"type:uuid;index;not null"`
	ImageType string    `gorm:"not null"`
	ImageURL  string    `gorm:"not null"`
}

// StoreTag represents a tag associated with a store.
type StoreTag struct {
	BaseModel
	TagName string `gorm:"not null;unique"`
}

// StoreHour represents the operating hours of a store.
type StoreHour struct {
	BaseModel
	StoreID   uuid.UUID     `gorm:"type:uuid;index;not null"`
	DayOfWeek DayOfWeekEnum `gorm:"not null"`
	Open      string
	Close     string
}
