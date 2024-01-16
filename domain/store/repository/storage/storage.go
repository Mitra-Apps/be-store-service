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
	validFileExtensions := []string{".png", ".jpg", ".jpeg", ".svg"}
	var fileExtension string

	for _, ext := range validFileExtensions {
		if mime.TypeByExtension(ext) == fileType {
			fileExtension = ext
			break
		}
	}

	if fileExtension == "" {
		return "", fmt.Errorf("file extension must be png, jpg, jpeg, or svg")
	}

	if len(decodedImage) > 4*1024*1024 {
		return "", fmt.Errorf("file size must be less than 4mb")
	}

	filename := uuid.New().String() + fileExtension
	objectName := fmt.Sprintf("stores/%s/%s", userID, filename)

	_, err = s.client.PutObject(ctx, s.bucket, objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	minioEndpoint := os.Getenv("STORAGE_ENDPOINT")
	return fmt.Sprintf("http://%s/%s/%s", minioEndpoint, s.bucket, objectName), nil
}
