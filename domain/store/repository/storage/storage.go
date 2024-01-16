package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/Mitra-Apps/be-store-service/domain/store/repository"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type storage struct {
	client *minio.Client
	bucket string
}

func New() repository.Storage {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	bucketName := os.Getenv("STORAGE_BUCKET_NAME")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		logrus.Errorln("Failed to connect to MinIO:", err)
	}

	return &storage{
		client: minioClient,
		bucket: bucketName,
	}
}

func (s *storage) UploadImage(ctx context.Context, image, userID string) (string, error) {
	b64data := image[strings.IndexByte(image, ',')+1:]
	decodedImage, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	fileType := http.DetectContentType(decodedImage)
	fileExtension, err := mime.ExtensionsByType(fileType)
	if err != nil || len(fileExtension) == 0 {
		return "", fmt.Errorf("failed to get file extension: %w", err)
	}

	// file extension must be png, jpg, jpeg, or svg
	if fileExtension[0] != ".png" && fileExtension[0] != ".jpg" && fileExtension[0] != ".jpeg" && fileExtension[0] != ".svg" {
		return "", fmt.Errorf("file extension must be png, jpg, jpeg, or svg")
	}

	// max file size is 4mb
	if len(decodedImage) > 4*1024*1024 {
		return "", fmt.Errorf("file size must be less than 4mb")
	}

	filename := uuid.New().String() + fileExtension[0]
	objectName := fmt.Sprintf("stores/%s/%s", userID, filename)

	_, err = s.client.PutObject(ctx, s.bucket, objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), s.bucket, objectName), nil
}

func getImageType(image string) string {
	if strings.HasPrefix(image, "data:image/jpeg;base64") {
		return "data:image/jpeg;base64"
	} else if strings.HasPrefix(image, "data:image/png;base64") {
		return "data:image/png;base64"
	} else if strings.HasPrefix(image, "data:image/gif;base64") {
		return "data:image/gif;base64"
	} else {
		return ""
	}
}

func removeImageTypePrefix(image string, imageType string) string {
	return strings.TrimPrefix(image, imageType)
}
