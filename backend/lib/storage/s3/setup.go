package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ErrFileTooBig struct {
	fileSize    int64
	fileName    string
	maxFileSize int64
}

func (e ErrFileTooBig) Error() string {
	return fmt.Sprintf("file %s is too big. Max size: %d, file size: %d", e.fileName, e.maxFileSize, e.fileSize)
}

type ErrSpaceLimitExceeded struct {
	maxUsage     int64
	currentUsage int64
}

func (e ErrSpaceLimitExceeded) Error() string {
	return fmt.Sprintf("max storage limit exceeded. Space limit: %d, currently used: %d", e.maxUsage, e.currentUsage)
}

// S3Storage allows to interact with S3 to store files
type S3Storage struct {
	s3                *s3.Client
	config            aws.Config
	bucket            string
	maxBucketSize     int64
	maxFileSize       int64
	presignExpiration time.Duration
}

// NewS3Storage creates a new S3 storage handler
func NewS3Storage(ctx context.Context, accessKeyID, secretKey, bucket string, options ...S3Option) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	storage := &S3Storage{
		s3:                s3.NewFromConfig(cfg),
		config:            cfg,
		bucket:            bucket,
		presignExpiration: 30,
	}

	for _, opt := range options {
		opt(storage)
	}

	return storage, nil
}

type S3Option func(*S3Storage)

func WithMaxBucketSize(size int64) S3Option {
	return func(s *S3Storage) {
		s.maxBucketSize = size
	}
}

func WithMaxFileSize(size int64) S3Option {
	return func(s *S3Storage) {
		s.maxFileSize = size
	}
}

func WithRegion(region string) S3Option {
	return func(s *S3Storage) {
		s.s3 = s3.NewFromConfig(s.config, func(o *s3.Options) {
			o.Region = region
		})
	}
}

func WithEndpoint(endpoint string) S3Option {
	return func(s *S3Storage) {
		s.s3 = s3.NewFromConfig(s.config, func(o *s3.Options) {
			o.BaseEndpoint = &endpoint
		})
	}
}
