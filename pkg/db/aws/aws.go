package aws

import (
	"context"
	"fmt"
	"github.com/AleksK1NG/api-mc/pkg/logger"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

// Minio AWS S3 Client interface
type AWSClient interface {
	FileUpload(ctx context.Context, input UploadInput) (minio.UploadInfo, error)
}

// AWS Upload Input
type UploadInput struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
	BucketName  string
}

// Minio AWS S3 Client
type AWSS3Client struct {
	client *minio.Client
}

// Minio AWS S3 Client constructor
func NewAWSClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) (*AWSS3Client, error) {

	// Initialize minio client object.
	minioClient, err := minio.New("play.min.io", &minio.Options{
		Creds:  credentials.NewStaticV4("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	logger.Infof("minioClient: %#v", minioClient)

	return &AWSS3Client{client: minioClient}, nil
}

// AWS file upload
func (aws *AWSS3Client) FileUpload(ctx context.Context, input UploadInput) (minio.UploadInfo, error) {
	options := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	uploadInfo, err := aws.client.PutObject(ctx, input.BucketName, input.Name, input.File, input.Size, options)
	if err != nil {
		logger.Errorf("FileUpload ", err)
		return uploadInfo, err
	}

	logger.Infof("AWS FileUpload: %#v", uploadInfo)
	return uploadInfo, err
}

func (aws *AWSS3Client) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", fileName, uid)
}
