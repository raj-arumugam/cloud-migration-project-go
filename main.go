package main

import (
	"log"

	"cloud-migration/internal/cloud"
	"cloud-migration/internal/cloud/aws"
	"cloud-migration/internal/cloud/google"
	"cloud-migration/internal/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var awsService cloud.PhotoService = aws.NewAWSPhotoService(&config.AWSConfig)
	var googleService cloud.PhotoService = google.NewGoogleDriveService(&config.GoogleConfig)

	// Connect to both services
	if err := awsService.Connect(); err != nil {
		log.Fatalf("Failed to connect to AWS: %v", err)
	}

	if err := googleService.Connect(); err != nil {
		log.Fatalf("Failed to connect to Google Drive: %v", err)
	}

	// List photos from AWS
	photos, err := awsService.ListPhotos()
	if err != nil {
		log.Fatalf("Failed to list photos from AWS: %v", err)
	}

	// Migrate each photo
	for _, photo := range photos {
		// Download from AWS
		data, err := awsService.DownloadPhoto(photo)
		if err != nil {
			log.Printf("Failed to download photo %s: %v", photo.Name, err)
			continue
		}

		// Upload to Google Drive
		if err := googleService.UploadPhoto(photo, data); err != nil {
			log.Printf("Failed to upload photo %s: %v", photo.Name, err)
			continue
		}

		log.Printf("Successfully migrated photo: %s", photo.Name)
	}

	log.Println("Migration completed!")
}
