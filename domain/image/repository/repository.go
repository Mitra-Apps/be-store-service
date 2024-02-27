package repository

import (
	"context"

	utilityPb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"github.com/google/uuid"
)

type ImageRepository interface {
	UploadImage(ctx context.Context, imageBase64Str, groupName, userID string) (*uuid.UUID, error)
	GetImagesByIds(ctx context.Context, ids []string) ([]*utilityPb.Image, error)
}
