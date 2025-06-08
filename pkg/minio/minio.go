package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/saufiroja/fin-ai/config"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type MinioManager interface {
	UploadFileFromMultipart(bucketName, objectName string, fileHeader *multipart.FileHeader) error
	FileExists(bucketName, objectName string) (bool, error)
	ReadAndEncodeFile(bucketName, objectName string) ([]byte, error)
	BucketExists(bucketName string) (bool, error)
	CreateBucket(bucketName string) error
}

type MinioClient struct {
	client *minio.Client
	conf   *config.AppConfig
}

var (
	instance *MinioClient
	once     sync.Once
)

func NewMinioClient(conf *config.AppConfig, log logging.Logger) MinioManager {
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

func (m *MinioClient) UploadFileFromMultipart(bucketName, objectName string, fileHeader *multipart.FileHeader) error {
	ctx := context.Background()

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Determine content type based on file extension
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".jpg") ||
		strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".png") {
		contentType = "image/png"
	}

	// Upload the file
	info, err := m.client.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("File uploaded successfully: %+v\n", info)
	return nil
}

func (m *MinioClient) FileExists(bucketName, objectName string) (bool, error) {
	ctx := context.Background()

	_, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil // File does not exist
		}
		return false, fmt.Errorf("error checking file existence: %w", err)
	}

	return true, nil // File exists
}

func (m *MinioClient) ReadAndEncodeFile(bucketName, objectName string) ([]byte, error) {
	ctx := context.Background()

	object, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting object: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("error reading object: %w", err)
	}

	return data, nil
}

func (m *MinioClient) BucketExists(bucketName string) (bool, error) {
	ctx := context.Background()

	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("error checking bucket existence: %w", err)
	}

	return exists, nil
}

func (m *MinioClient) CreateBucket(bucketName string) error {
	ctx := context.Background()

	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if exists {
		return nil // Bucket already exists
	}

	err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	log.Printf("Bucket %s created successfully\n", bucketName)
	return nil
}
