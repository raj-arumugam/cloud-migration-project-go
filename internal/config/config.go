package config

import (
	"cloud-migration/internal/cloud/aws"
	"cloud-migration/internal/cloud/google"
)

type Config struct {
	AWSConfig    aws.Config
	GoogleConfig google.Config
}

func LoadConfig() (*Config, error) {
	// TODO: Implement configuration loading from environment variables or config file
	return &Config{}, nil
}
