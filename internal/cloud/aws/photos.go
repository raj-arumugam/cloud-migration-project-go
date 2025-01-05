package aws

import (
	"context"
	"fmt"

	"cloud-migration/internal/cloud"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSPhotoService struct {
	config *Config
	client *s3.Client
}

func NewAWSPhotoService(config *Config) *AWSPhotoService {
	return &AWSPhotoService{
		config: config,
	}
}

func (s *AWSPhotoService) Connect() error {
	if s.config.AccessKeyID == "" || s.config.SecretAccessKey == "" {
		return fmt.Errorf("AWS credentials not provided")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	s.client = s3.NewFromConfig(cfg)
	return nil
}

func (s *AWSPhotoService) ListPhotos() ([]cloud.Photo, error) {
	if s.client == nil {
		return nil, fmt.Errorf("AWS client not initialized")
	}

	// This will now return an error if bucket is not configured
	if s.config.BucketName == "" {
		return nil, fmt.Errorf("AWS bucket name not configured")
	}

	return nil, fmt.Errorf("not yet implemented")
}

func (s *AWSPhotoService) DownloadPhoto(photo cloud.Photo) ([]byte, error) {
	// TODO: Use AWS SDK to download photo from the configured bucket
	return nil, nil
}

func (s *AWSPhotoService) UploadPhoto(photo cloud.Photo, data []byte) error {
	// TODO: Use AWS SDK to upload photo to the configured bucket
	return nil
}
