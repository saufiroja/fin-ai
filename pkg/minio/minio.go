package minio

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/saufiroja/fin-ai/config"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type MinioManager interface {
	UploadFile(bucketName, objectName string, filePath string) error
}

type MinioClient struct {
	client *minio.Client
	conf   *config.AppConfig
}

func NewMinioClient(conf *config.AppConfig, logging logging.Logger) *MinioClient {
	// Inisialisasi client MinIO
	client, err := minio.New(conf.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Minio.AccessKey, conf.Minio.SecretKey, ""),
		Secure: conf.Minio.UseSSL,
	})
	if err != nil {
		logging.LogError(fmt.Sprintf("Failed to create MinIO client: %v", err))
		return nil
	}

	logging.LogInfo(fmt.Sprintf("MinIO client connected to: %s", conf.Minio.Endpoint))
	return &MinioClient{
		client: client,
		conf:   conf,
	}
}

func (m *MinioClient) UploadFile(bucketName, objectName string, filePath string) error {
	ctx := context.Background()

	// Cek apakah bucket sudah ada
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if !exists {
		// Jika tidak ada, buat bucket
		if err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully\n", bucketName)
	}

	// Upload file
	contentType := "application/octet-stream" // default, bisa disesuaikan
	info, err := m.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("File uploaded successfully: %+v\n", info)
	return nil
}
