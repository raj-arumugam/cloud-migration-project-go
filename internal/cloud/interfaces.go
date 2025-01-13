package cloud

import (
	"context"
	"time"
)

type PhotoService interface {
	Connect(ctx context.Context) error
	ListPhotos(ctx context.Context) ([]Photo, error)
	DownloadPhoto(ctx context.Context, photo Photo) ([]byte, error)
	UploadPhoto(ctx context.Context, photo Photo, data []byte) error
}

type Photo struct {
	ID        string
	Name      string
	Path      string                 // Optional: local file path
	MimeType  string                 // File content type
	Size      int64                  // File size in bytes
	CreatedAt *time.Time             // Creation timestamp
	Metadata  map[string]interface{} // Optional metadata
}
