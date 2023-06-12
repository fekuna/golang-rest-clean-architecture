// package repository

// import (
// 	"github.com/fekuna/go-rest-clean-architecture/internal/auth"
// 	"github.com/fekuna/go-rest-clean-architecture/pkg/db/minioS3"
// )

// // Auth Minio S3 repository
// type authMinioRepository struct {
// 	minioConfig *minioS3.MinioConfig
// }

// // Auth Minio S3 Repository constructor
// func NewAuthMinioRepository(minioConfig *minioS3.MinioConfig) auth.MinioRepository {
// 	return &authMinioRepository{minioConfig: minioConfig}
// }

// Upload file to Minio
// func (r *authMinioRepository) PutObject(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error) {
// 	// TODO: Tracing

// 	uploadInfo, err := r.client.PutObject(ctx, input.BucketName, r.generateFileName(input.Name), input.File, input.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
// 	if err != nil {
// 		fmt.Println("Mashok", err)
// 		return nil, errors.Wrap(err, "authAWSRepository.FileUpload.PutObject")
// 	}
// 	return &uploadInfo, err
// }

// // Download file from minio
// func (r *authMinioRepository) GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error) {
// 	// TODO: Tracing

// 	object, err := r.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "authAWSRepository.FileDownload.GetObject")
// 	}

// 	return object, nil
// }

// // Delete file from Minio
// func (r *authMinioRepository) RemoveObject(ctx context.Context, bucket string, fileName string) error {
// 	// TODO: Tracing

// 	if err := r.client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{}); err != nil {
// 		return errors.Wrap(err, "authAWSRepository.RemoveObject")
// 	}

// 	return nil
// }

// // GetUrl of object from Minio
// func (r *authMinioRepository) GetObjectUrl(ctx context.Context, bucket string, fileName string, expires time.Duration) (*url.URL, error) {
// 	if expires == 0 {
// 		expires = time.Second * 604800
// 	}
// 	// reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

// 	objectUrl, err := r.client.PresignedGetObject(ctx, bucket, fileName, expires, url.Values{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return objectUrl, nil
// }

//	func (r *authMinioRepository) generateFileName(fileName string) string {
//		uid := uuid.New().String()
//		return fmt.Sprintf("%s-%s", uid, fileName)
//	}
package repository
