package repository

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// Auth AWS S3 repository
type authAWSRepository struct {
	client *minio.Client
}

// Auth AWS S3 repository constructor
func NewAuthAWSRepository(awsClient *minio.Client) auth.AWSRepository {
	return &authAWSRepository{client: awsClient}
}

// Upload file to AWS
func (aws *authAWSRepository) FileUpload(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error) {
	options := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	uploadInfo, err := aws.client.PutObject(ctx, input.BucketName, aws.generateFileName(input.Name), input.File, input.Size, options)
	if err != nil {
		return nil, err
	}

	return &uploadInfo, err
}

func (aws *authAWSRepository) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uid, fileName)
}
