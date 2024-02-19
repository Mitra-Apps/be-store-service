package entity

import (
	"regexp"

	"github.com/Mitra-Apps/be-store-service/domain/base_model"
	prodEntity "github.com/Mitra-Apps/be-store-service/domain/product/entity"
	pb "github.com/Mitra-Apps/be-store-service/domain/proto/store"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Store represents a store model.
type Store struct {
	base_model.BaseModel
	UserID           uuid.UUID `gorm:"type:uuid;not null;unique"`
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
	Tags             []*StoreTag           `gorm:"many2many:store_store_tags"`
	Hours            []*StoreHour          `gorm:"foreignKey:StoreID"`
	Images           []*StoreImage         `gorm:"foreignKey:StoreID"`
	Products         []*prodEntity.Product `gorm:"foreignKey:StoreID"`
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
			return status.Errorf(codes.InvalidArgument, "invalid store id")
		}
		s.ID = id
	}

	// convert user id to uuid
	var userID uuid.UUID
	if store.UserId != "" {
		id, err := uuid.Parse(store.UserId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid user id")
		}
		userID = id
	}

	// check user email is valid email address
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(store.Email) {
		return status.Errorf(codes.InvalidArgument, "invalid email address")
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
		s.Tags = append(s.Tags, storeTag)
	}

	for _, hour := range store.Hours {
		storeHour := &StoreHour{}
		if err := storeHour.FromProto(hour); err != nil {
			return err
		}
		s.Hours = append(s.Hours, storeHour)
	}

	for _, image := range store.Images {
		storeImage := &StoreImage{}
		if err := storeImage.FromProto(image); err != nil {
			return err
		}
		s.Images = append(s.Images, storeImage)
	}

	return nil
}

// StoreImage represents an image associated with a store.
type StoreImage struct {
	base_model.BaseModel
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
	base_model.BaseModel
	TagName string
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
	base_model.BaseModel
	StoreID   uuid.UUID `gorm:"type:uuid;index;not null"`
	DayOfWeek string    `gorm:"not null"`
	Open      string
	Close     string
	Is24Hr    bool
	IsOpen    bool `gorm:"type:bool;not null;default:true"`
}

func (s *StoreHour) ToProto() *pb.StoreHour {
	var dayOfWeek int32

	switch s.DayOfWeek {
	case "MONDAY":
		dayOfWeek = 0
	case "TUESDAY":
		dayOfWeek = 1
	case "WEDNESDAY":
		dayOfWeek = 2
	case "THURSDAY":
		dayOfWeek = 3
	case "FRIDAY":
		dayOfWeek = 4
	case "SATURDAY":
		dayOfWeek = 5
	case "SUNDAY":
		dayOfWeek = 6
	}

	return &pb.StoreHour{
		Id:        s.ID.String(),
		StoreId:   s.StoreID.String(),
		DayOfWeek: dayOfWeek,
		Open:      s.Open,
		Close:     s.Close,
		Is24Hours: s.Is24Hr,
		IsOpen:    s.IsOpen,
	}
}

func (s *StoreHour) FromProto(storeHour *pb.StoreHour) error {
	switch storeHour.DayOfWeek {
	case 0:
		s.DayOfWeek = "MONDAY"
	case 1:
		s.DayOfWeek = "TUESDAY"
	case 2:
		s.DayOfWeek = "WEDNESDAY"
	case 3:
		s.DayOfWeek = "THURSDAY"
	case 4:
		s.DayOfWeek = "FRIDAY"
	case 5:
		s.DayOfWeek = "SATURDAY"
	case 6:
		s.DayOfWeek = "SUNDAY"
	}

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

	s.Open = storeHour.Open
	s.Close = storeHour.Close
	s.Is24Hr = storeHour.Is24Hours
	s.IsOpen = storeHour.IsOpen

	return nil
}
