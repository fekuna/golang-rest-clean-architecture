//go:generate mockgen -source minio_repository.go -destination mock/minio_repository_mock.go -package mock
package auth

// Minio S3 interface
// type MinioRepository interface {
// 	PutObject(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error)
// 	GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error)
// 	RemoveObject(ctx context.Context, bucket string, fileName string) error
// 	GetObjectUrl(ctx context.Context, bucket string, fileName string, expires time.Duration) (*url.URL, error)
// }
