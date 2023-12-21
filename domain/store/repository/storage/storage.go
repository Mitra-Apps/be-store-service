package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"

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
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Fatalln("Error creating MinIO client:", err)
	}

	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		logrus.Fatalln("Error checking if bucket exists:", err)
	}
	if !exists {
		err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Fatalln("Error creating bucket:", err)
		}
	}

	return &storage{
		client: minioClient,
		bucket: bucketName,
	}
}

func (s *storage) UploadImage(ctx context.Context, image string) (string, error) {
	objectName := uuid.New().String()

	decodedImage, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		fmt.Println("Error decoding base64 image:", err)
		return "", err
	}

	fileType := http.DetectContentType(decodedImage)

	fileExtension, err := mime.ExtensionsByType(fileType)
	if err != nil || len(fileExtension) == 0 {
		fmt.Println("Error determining file extension:", err)
		return "", err
	}

	objectName += fileExtension[0]

	_, err = s.client.PutObject(ctx, s.bucket, objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("Error uploading image:", err)
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), s.bucket, objectName), nil
}
