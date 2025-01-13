package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"cloud-migration/internal/cloud"
	"cloud-migration/internal/cloud/aws"
	cloudgoogle "cloud-migration/internal/cloud/google"
	"cloud-migration/internal/config"
	"cloud-migration/internal/migrator"
)

func setupLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func setupServices(cfg *config.Config) (cloud.PhotoService, cloud.PhotoService, error) {
	awsService := aws.NewAWSPhotoService(&aws.Config{
		Region:          cfg.AWSConfig.Region,
		Bucket:          cfg.AWSConfig.Bucket,
		AccessKeyID:     cfg.AWSConfig.AccessKeyID,
		SecretAccessKey: cfg.AWSConfig.SecretAccessKey,
		RateLimit:       cfg.AWSConfig.RateLimit,
	})

	googleService := cloudgoogle.NewGoogleDriveService(&cloudgoogle.GoogleDriveConfig{
		ClientID:     cfg.GoogleConfig.ClientID,
		ClientSecret: cfg.GoogleConfig.ClientSecret,
		TokenPath:    cfg.GoogleConfig.TokenPath,
		RateLimit:    cfg.GoogleConfig.RateLimit,
	})

	return awsService, googleService, nil
}

func main() {
	// Create root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup logger
	logger, err := setupLogger()
	if err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config",
			zap.Error(err),
		)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Invalid configuration",
			zap.Error(err),
		)
	}

	// Setup services
	sourceService, destService, err := setupServices(cfg)
	if err != nil {
		logger.Fatal("Failed to setup services",
			zap.Error(err),
		)
	}

	// Create migrator
	photoMigrator := migrator.NewPhotoMigrator(sourceService, destService, logger)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal",
			zap.String("signal", sig.String()),
		)
		cancel()
	}()

	// Log starting migration
	logger.Info("Starting migration",
		zap.String("source", "AWS"),
		zap.String("destination", "Google Drive"),
		zap.String("bucket", cfg.AWSConfig.Bucket))

	// Run migration
	if err := photoMigrator.MigratePhotos(ctx); err != nil {
		logger.Error("Migration failed",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Info("Migration completed successfully")
}
