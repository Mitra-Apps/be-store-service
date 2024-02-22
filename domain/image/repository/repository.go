package repository

import (
	"context"

	"github.com/google/uuid"
)

type ImageRepository interface {
	UploadImage(ctx context.Context, imageBase64Str, groupName, userID string) (*uuid.UUID, error)
}
