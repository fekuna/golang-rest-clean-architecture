package minioS3

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type MinioConfig struct {
	client *minio.Client
}

// Minio S3 Client constructor
func NewMinioS3Client(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) (*MinioConfig, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioConfig{
		client: minioClient,
	}, nil
}

// Bucket exists checker
func (r *MinioConfig) IsBucketExists(ctx context.Context, bucketName string) (bool, error) {
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		return true, err
	}

	return exists, nil
}

// Create bucket on minio
func (r *MinioConfig) CreateBucket(ctx context.Context, bucketName string) error {
	isBucketExists, err := r.IsBucketExists(ctx, bucketName)
	if err != nil || isBucketExists {
		return err
	}

	return r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
}

// Upload file to Minio
func (r *MinioConfig) PutObject(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error) {
	// TODO: Tracing

	uploadInfo, err := r.client.PutObject(ctx, input.BucketName, r.generateFileName(input.Name), input.File, input.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, errors.Wrap(err, "authAWSRepository.FileUpload.PutObject")
	}
	return &uploadInfo, err
}

// Download file from minio
func (r *MinioConfig) GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error) {
	// TODO: Tracing

	object, err := r.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "authAWSRepository.FileDownload.GetObject")
	}

	return object, nil
}

// Delete file from Minio
func (r *MinioConfig) RemoveObject(ctx context.Context, bucket string, fileName string) error {
	// TODO: Tracing

	if err := r.client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{}); err != nil {
		return errors.Wrap(err, "authAWSRepository.RemoveObject")
	}

	return nil
}

// GetUrl of object from Minio
func (r *MinioConfig) GetObjectUrl(ctx context.Context, bucket string, fileName string, expires time.Duration) (*url.URL, error) {
	if expires == 0 {
		expires = time.Second * 604800
	}
	// reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

	objectUrl, err := r.client.PresignedGetObject(ctx, bucket, fileName, expires, url.Values{})
	if err != nil {
		return nil, err
	}

	return objectUrl, nil
}

func (r *MinioConfig) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uid, fileName)
}
