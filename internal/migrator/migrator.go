package migrator

import (
	"cloud-migration/internal/cloud"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap" // Add this import
)

type PhotoMigrator struct {
	sourceService      cloud.PhotoService
	destinationService cloud.PhotoService
	logger             *zap.Logger
	metrics            *Metrics
}

type Metrics struct {
	PhotosMigrated   int64
	BytesTransferred int64
}

func NewPhotoMigrator(source, destination cloud.PhotoService, logger *zap.Logger) *PhotoMigrator {
	return &PhotoMigrator{
		sourceService:      source,
		destinationService: destination,
		logger:             logger,
		metrics:            &Metrics{},
	}
}

func (m *PhotoMigrator) Connect(ctx context.Context) error {
	if err := m.sourceService.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to source service: %w", err)
	}

	if err := m.destinationService.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to destination service: %w", err)
	}

	return nil
}

func (m *PhotoMigrator) MigratePhotos(ctx context.Context) error {
	var errs []error
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("Panic recovered",
				zap.Any("panic", r),
			)
		}
	}()

	photos, err := m.sourceService.ListPhotos(ctx)
	if err != nil {
		return fmt.Errorf("failed to list photos: %w", err)
	}

	for _, photo := range photos {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := m.migratePhoto(ctx, photo); err != nil {
				errs = append(errs, fmt.Errorf("failed to migrate photo %s: %w", photo.Name, err))
				continue // Continue with next photo instead of stopping
			}
			m.logger.Info("successfully migrated photo",
				zap.String("photo", photo.Name))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("migration completed with %d errors: %v", len(errs), errs)
	}
	return nil
}

func (m *PhotoMigrator) migratePhoto(ctx context.Context, photo cloud.Photo) error {
	backoff := time.Second
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		err := m.attemptPhotoMigration(ctx, photo)
		if err == nil {
			return nil
		}

		m.logger.Error("Migration attempt failed",
			zap.String("photo", photo.Name),
			zap.Error(err),
			zap.Int("attempt", i+1))

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			backoff *= 2
		}
	}
	return fmt.Errorf("failed to migrate photo after %d attempts", maxRetries)
}

func (m *PhotoMigrator) attemptPhotoMigration(ctx context.Context, photo cloud.Photo) error {
	if photo.Name == "" {
		return fmt.Errorf("invalid photo: name is empty")
	}

	var data []byte
	var err error

	// Retry download up to 3 times
	for retries := 0; retries < 3; retries++ {
		data, err = m.sourceService.DownloadPhoto(ctx, photo)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(retries+1))
	}
	if err != nil {
		return fmt.Errorf("failed to download photo after retries: %w", err)
	}

	// Retry upload up to 3 times
	for retries := 0; retries < 3; retries++ {
		err = m.destinationService.UploadPhoto(ctx, photo, data)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(retries+1))
	}
	if err != nil {
		return fmt.Errorf("failed to upload photo after retries: %w", err)
	}

	atomic.AddInt64(&m.metrics.PhotosMigrated, 1)
	atomic.AddInt64(&m.metrics.BytesTransferred, int64(len(data)))

	return nil
}
