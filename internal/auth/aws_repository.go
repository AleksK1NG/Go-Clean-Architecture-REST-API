package auth

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/minio/minio-go/v7"
)

// Minio AWS S3 interface
type AWSRepository interface {
	FileUpload(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error)
	FileDownload(ctx context.Context, bucket string, fileName string) (*minio.Object, error)
}
