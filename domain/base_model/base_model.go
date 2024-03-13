package base_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;null"`
	UpdatedBy uuid.UUID       `gorm:"type:uuid;null"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz;null"`
	DeletedBy uuid.UUID      `gorm:"type:uuid;null"`
}

type BaseMasterDataModel struct {
	ID        int64          `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;null"`
	UpdatedBy uuid.UUID      `gorm:"type:uuid;null"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz;null"`
	DeletedBy uuid.UUID      `gorm:"type:uuid;null"`
}
