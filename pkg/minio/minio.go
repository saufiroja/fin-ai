package minio

import (
	"context"
	"fmt"
	"log"
	"sync"

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

var (
	instance *MinioClient
	once     sync.Once
)

func NewMinioClient(conf *config.AppConfig, log logging.Logger) *MinioClient {
	once.Do(func() {
		client, err := minio.New(conf.Minio.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.Minio.AccessKey, conf.Minio.SecretKey, ""),
			Secure: conf.Minio.UseSSL,
		})
		if err != nil {
			log.LogError(fmt.Sprintf("Failed to create MinIO client: %v", err))
			return
		}

		log.LogInfo(fmt.Sprintf("MinIO client connected to: %s", conf.Minio.Endpoint))
		instance = &MinioClient{
			client: client,
			conf:   conf,
		}
	})

	return instance
}

func (m *MinioClient) UploadFile(bucketName, objectName string, filePath string) error {
	ctx := context.Background()

	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if !exists {
		if err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully\n", bucketName)
	}

	contentType := "application/octet-stream"
	info, err := m.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("File uploaded successfully: %+v\n", info)
	return nil
}
