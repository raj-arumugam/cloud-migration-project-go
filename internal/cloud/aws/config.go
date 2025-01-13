package aws

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	RateLimit       float64
}
