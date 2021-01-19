package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/api-mc/internal/auth"
	"github.com/AleksK1NG/api-mc/internal/models"
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
func (aws *authAWSRepository) PutObject(ctx context.Context, input models.UploadInput) (*minio.UploadInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authAWSRepository.PutObject")
	defer span.Finish()

	options := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	uploadInfo, err := aws.client.PutObject(ctx, input.BucketName, aws.generateFileName(input.Name), input.File, input.Size, options)
	if err != nil {
		return nil, errors.Wrap(err, "authAWSRepository.FileUpload.PutObject")
	}

	return &uploadInfo, err
}

// Download file from AWS
func (aws *authAWSRepository) GetObject(ctx context.Context, bucket string, fileName string) (*minio.Object, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authAWSRepository.GetObject")
	defer span.Finish()

	object, err := aws.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "authAWSRepository.FileDownload.GetObject")
	}
	return object, nil
}

// Delete file from AWS
func (aws *authAWSRepository) RemoveObject(ctx context.Context, bucket string, fileName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authAWSRepository.RemoveObject")
	defer span.Finish()

	if err := aws.client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{}); err != nil {
		return errors.Wrap(err, "authAWSRepository.RemoveObject")
	}
	return nil
}

func (aws *authAWSRepository) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uid, fileName)
}
