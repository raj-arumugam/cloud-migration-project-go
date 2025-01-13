package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AWSConfig struct {
		Region          string
		Bucket          string
		RateLimit       float64
		AccessKeyID     string
		SecretAccessKey string
	}
	GoogleConfig struct {
		ClientID     string
		ClientSecret string
		TokenPath    string
		RateLimit    float64
	}
	RetryAttempts int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// AWS Config
	cfg.AWSConfig.Region = os.Getenv("AWS_REGION")
	cfg.AWSConfig.Bucket = os.Getenv("AWS_BUCKET_NAME")
	cfg.AWSConfig.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	cfg.AWSConfig.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	rateLimitStr := os.Getenv("AWS_RATE_LIMIT")
	if rateLimitStr != "" {
		rateLimit, err := strconv.ParseFloat(rateLimitStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid AWS_RATE_LIMIT: %w", err)
		}
		cfg.AWSConfig.RateLimit = rateLimit
	} else {
		cfg.AWSConfig.RateLimit = 10 // default
	}

	// Google Config
	cfg.GoogleConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.GoogleConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.GoogleConfig.TokenPath = os.Getenv("GOOGLE_TOKEN_PATH")

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.AWSConfig.Region == "" {
		return fmt.Errorf("AWS region is required")
	}
	if c.AWSConfig.Bucket == "" {
		return fmt.Errorf("AWS bucket is required")
	}
	if c.AWSConfig.AccessKeyID == "" {
		return fmt.Errorf("AWS access key ID is required")
	}
	if c.AWSConfig.SecretAccessKey == "" {
		return fmt.Errorf("AWS secret access key is required")
	}
	if c.GoogleConfig.ClientID == "" {
		return fmt.Errorf("Google client ID is required")
	}
	if c.GoogleConfig.ClientSecret == "" {
		return fmt.Errorf("Google client secret is required")
	}
	if c.GoogleConfig.TokenPath == "" {
		return fmt.Errorf("Google token path is required")
	}
	if c.AWSConfig.RateLimit <= 0 {
		return fmt.Errorf("AWS rate limit must be positive")
	}
	if c.GoogleConfig.RateLimit <= 0 {
		return fmt.Errorf("Google rate limit must be positive")
	}
	if c.RetryAttempts <= 0 {
		return fmt.Errorf("retry attempts must be positive")
	}
	return nil
}
