package google

import (
	"cloud-migration/internal/cloud"
)

type GoogleDriveService struct {
	config *Config
}

// NewGoogleDriveService creates a new Google Drive service instance
func NewGoogleDriveService(config *Config) *GoogleDriveService {
	return &GoogleDriveService{
		config: config,
	}
}

func (s *GoogleDriveService) Connect() error {
	// TODO: Initialize Google Drive client using s.config.ClientID, s.config.ClientSecret, etc.
	return nil
}

func (s *GoogleDriveService) ListPhotos() ([]cloud.Photo, error) {
	// TODO: Use Google Drive API to list photos
	return nil, nil
}

func (s *GoogleDriveService) DownloadPhoto(photo cloud.Photo) ([]byte, error) {
	// TODO: Use Google Drive API to download photo
	return nil, nil
}

func (s *GoogleDriveService) UploadPhoto(photo cloud.Photo, data []byte) error {
	// TODO: Use Google Drive API to upload photo
	return nil
}
