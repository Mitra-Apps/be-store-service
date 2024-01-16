package entity

import (
	"time"

	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
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
	UserID           uuid.UUID `gorm:"type:uuid;not null"`
	StoreName        string    `gorm:"not null;unique"`
	StoreDescription string    `gorm:"not null,type:text"`
	Address          string    `gorm:"not null,type:text"`
	City             string    `gorm:"not null"`
	State            string    `gorm:"not null"`
	ZipCode          string    `gorm:"not null"`
	Phone            string    `gorm:"not null"`
	Email            string    `gorm:"not null;unique"`
	Website          string
	LocationLat      float64
	LocationLng      float64
	Status           string
	IsActive         bool
	Tags             []StoreTag   `gorm:"many2many:store_store_tags"`
	Hours            []StoreHour  `gorm:"foreignKey:StoreID"`
	Images           []StoreImage `gorm:"foreignKey:StoreID"`
}

func (s *Store) ToProto() *pb.Store {
	tags := []*pb.StoreTag{}
	for _, tag := range s.Tags {
		tags = append(tags, tag.ToProto())
	}

	hours := []*pb.StoreHour{}
	for _, hour := range s.Hours {
		hours = append(hours, hour.ToProto())
	}

	images := []*pb.StoreImage{}
	for _, image := range s.Images {
		images = append(images, image.ToProto())
	}

	return &pb.Store{
		Id:               s.ID.String(),
		UserId:           s.UserID.String(),
		StoreName:        s.StoreName,
		StoreDescription: s.StoreDescription,
		Address:          s.Address,
		City:             s.City,
		State:            s.State,
		ZipCode:          s.ZipCode,
		Phone:            s.Phone,
		Email:            s.Email,
		Website:          s.Website,
		LocationLat:      s.LocationLat,
		LocationLng:      s.LocationLng,
		Status:           s.Status,
		IsActive:         s.IsActive,
		Tags:             tags,
		Hours:            hours,
		Images:           images,
	}
}

func (s *Store) FromProto(store *pb.Store) error {
	if store.Id != "" {
		id, err := uuid.Parse(store.Id)
		if err != nil {
			return err
		}
		s.ID = id
	}

	// convert user id to uuid
	userID, err := uuid.Parse(store.UserId)
	if err != nil {
		return err
	}

	s.UserID = userID
	s.StoreName = store.StoreName
	s.StoreDescription = store.StoreDescription
	s.Address = store.Address
	s.City = store.City
	s.State = store.State
	s.ZipCode = store.ZipCode
	s.Phone = store.Phone
	s.Email = store.Email
	s.Website = store.Website
	s.LocationLat = store.LocationLat
	s.LocationLng = store.LocationLng
	s.Status = store.Status
	s.IsActive = store.IsActive

	for _, tag := range store.Tags {
		storeTag := &StoreTag{}
		if err := storeTag.FromProto(tag); err != nil {
			return err
		}
		s.Tags = append(s.Tags, *storeTag)
	}

	for _, hour := range store.Hours {
		storeHour := &StoreHour{}
		if err := storeHour.FromProto(hour); err != nil {
			return err
		}
		s.Hours = append(s.Hours, *storeHour)
	}

	for _, image := range store.Images {
		storeImage := &StoreImage{}
		if err := storeImage.FromProto(image); err != nil {
			return err
		}
		s.Images = append(s.Images, *storeImage)
	}

	return nil
}

// StoreImage represents an image associated with a store.
type StoreImage struct {
	BaseModel
	StoreID     uuid.UUID `gorm:"type:uuid;index;not null"`
	ImageType   string    `gorm:"not null"`
	ImageURL    string    `gorm:"not null"`
	ImageBase64 string    `gorm:"-"`
}

func (s *StoreImage) ToProto() *pb.StoreImage {
	return &pb.StoreImage{
		Id:          s.ID.String(),
		StoreId:     s.StoreID.String(),
		ImageType:   s.ImageType,
		ImageUrl:    s.ImageURL,
		ImageBase64: "",
	}
}

func (s *StoreImage) FromProto(storeImage *pb.StoreImage) error {
	if storeImage.Id != "" {
		id, err := uuid.Parse(storeImage.Id)
		if err != nil {
			return err
		}
		s.ID = id
	}

	if storeImage.StoreId != "" {
		storeID, err := uuid.Parse(storeImage.StoreId)
		if err != nil {
			return err
		}
		s.StoreID = storeID
	}

	s.ImageType = storeImage.ImageType
	s.ImageURL = storeImage.ImageUrl
	s.ImageBase64 = storeImage.ImageBase64

	return nil
}

// StoreTag represents a tag associated with a store.
type StoreTag struct {
	BaseModel
	TagName string `gorm:"not null;unique"`
}

func (s *StoreTag) ToProto() *pb.StoreTag {
	return &pb.StoreTag{
		Id:      s.ID.String(),
		TagName: s.TagName,
	}
}

func (s *StoreTag) FromProto(storeTag *pb.StoreTag) error {
	if storeTag.Id != "" {
		id, err := uuid.Parse(storeTag.Id)
		if err != nil {
			return err
		}
		s.ID = id
	}

	s.TagName = storeTag.TagName

	return nil
}

// StoreHour represents the operating hours of a store.
type StoreHour struct {
	BaseModel
	StoreID   uuid.UUID     `gorm:"type:uuid;index;not null"`
	DayOfWeek DayOfWeekEnum `gorm:"not null"`
	Open      string
	Close     string
}

func (s *StoreHour) ToProto() *pb.StoreHour {
	dayOfWeekEnum := pb.DayOfWeekEnum(pb.DayOfWeekEnum_value[string(s.DayOfWeek)])

	return &pb.StoreHour{
		Id:        s.ID.String(),
		StoreId:   s.StoreID.String(),
		DayOfWeek: dayOfWeekEnum,
		Open:      s.Open,
		Close:     s.Close,
	}
}

func (s *StoreHour) FromProto(storeHour *pb.StoreHour) error {
	dayOfWeekEnum := DayOfWeekEnum(pb.DayOfWeekEnum_name[int32(storeHour.DayOfWeek)])

	if storeHour.Id != "" {
		id, err := uuid.Parse(storeHour.Id)
		if err != nil {
			return err
		}
		s.ID = id
	}

	if storeHour.StoreId != "" {
		storeID, err := uuid.Parse(storeHour.StoreId)
		if err != nil {
			return err
		}
		s.StoreID = storeID
	}

	s.DayOfWeek = dayOfWeekEnum
	s.Open = storeHour.Open
	s.Close = storeHour.Close

	return nil
}
