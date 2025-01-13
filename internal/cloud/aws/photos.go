package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud-migration/internal/cloud"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/time/rate"
)

type AWSPhotoService struct {
	config      *Config
	client      *s3.Client
	bucket      string
	rateLimiter *rate.Limiter
}

func NewAWSPhotoService(cfg *Config) *AWSPhotoService {
	return &AWSPhotoService{
		config:      cfg,
		rateLimiter: rate.NewLimiter(rate.Limit(cfg.RateLimit), 5), // Allow burst of 5
	}
}

func (s *AWSPhotoService) Connect(ctx context.Context) error {
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

func (s *AWSPhotoService) ListPhotos(ctx context.Context) ([]cloud.Photo, error) {
	if s.client == nil {
		return nil, fmt.Errorf("AWS client not initialized")
	}

	if s.config.Bucket == "" {
		return nil, fmt.Errorf("AWS bucket name not configured")
	}

	var photos []cloud.Photo
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.config.Bucket,
	})

	for paginator.HasMorePages() {
		if err := s.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit error: %w", err)
		}

		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			photos = append(photos, cloud.Photo{
				ID:   *obj.Key,
				Name: *obj.Key,
				Path: *obj.Key,
			})
		}
	}

	return photos, nil
}

func (s *AWSPhotoService) DownloadPhoto(ctx context.Context, photo cloud.Photo) ([]byte, error) {
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	input := &s3.GetObjectInput{
		Bucket: &s.config.Bucket,
		Key:    &photo.Path,
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download photo: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *AWSPhotoService) UploadPhoto(ctx context.Context, photo cloud.Photo, data []byte) error {
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit error: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket: &s.config.Bucket,
		Key:    &photo.Path,
		Body:   bytes.NewReader(data),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload photo: %w", err)
	}

	return nil
}
