package google

import (
	"bytes"
	"cloud-migration/internal/cloud"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	"google.golang.org/api/drive/v3"
)

// GoogleDriveConfig holds configuration for the Google Drive service
type GoogleDriveConfig struct {
	RateLimit       float64 // requests per second
	BurstLimit      int     // burst size
	CredentialsPath string
	ClientID        string
	ClientSecret    string
	TokenPath       string
}

type GoogleDriveService struct {
	service     *drive.Service
	rateLimiter *rate.Limiter
	config      *GoogleDriveConfig
}

// NewGoogleDriveService creates a new Google Drive service instance
func NewGoogleDriveService(cfg *GoogleDriveConfig) *GoogleDriveService {
	return &GoogleDriveService{
		config:      cfg,
		rateLimiter: rate.NewLimiter(rate.Limit(cfg.RateLimit), 5), // Allow burst of 5
	}
}

func (g *GoogleDriveService) Connect(ctx context.Context) error {
	config := &oauth2.Config{
		ClientID:     g.config.ClientID,
		ClientSecret: g.config.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
	}

	tokenBytes, err := os.ReadFile(g.config.TokenPath)
	if err != nil {
		return fmt.Errorf("failed to read token file: %w", err)
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(tokenBytes, token); err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	client := config.Client(ctx, token)
	service, err := drive.New(client)
	if err != nil {
		return fmt.Errorf("failed to create Drive client: %w", err)
	}

	g.service = service
	return nil
}

func (g *GoogleDriveService) ListPhotos(ctx context.Context) ([]cloud.Photo, error) {
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	files, err := g.service.Files.List().
		Fields("files(id, name, mimeType, size, createdTime)").
		Q("mimeType contains 'image/'").
		PageSize(1000).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	if files == nil || len(files.Files) == 0 {
		return []cloud.Photo{}, nil
	}

	photos := make([]cloud.Photo, 0, len(files.Files))
	for _, file := range files.Files {
		if file == nil {
			continue
		}

		createdTime, err := time.Parse(time.RFC3339, file.CreatedTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created time for file %s: %w", file.Id, err)
		}

		photo := cloud.Photo{
			ID:        file.Id,
			Name:      file.Name,
			MimeType:  file.MimeType,
			Size:      file.Size,
			CreatedAt: &createdTime,
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (g *GoogleDriveService) DownloadPhoto(ctx context.Context, photo cloud.Photo) ([]byte, error) {
	if photo.ID == "" {
		return nil, fmt.Errorf("photo ID cannot be empty")
	}

	if err := g.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	resp, err := g.service.Files.Get(photo.ID).Context(ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (g *GoogleDriveService) UploadPhoto(ctx context.Context, photo cloud.Photo, data []byte) error {
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	file := &drive.File{
		Name:     photo.Name,
		MimeType: photo.MimeType,
	}

	_, err := g.service.Files.Create(file).
		Media(bytes.NewReader(data)).
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
