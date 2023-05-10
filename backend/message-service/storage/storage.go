package storage

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// StorageLayer defines functionality expected from Storage
type StorageLayer interface {
	GetPresignedPutRequest(key string) (string, error)
	DeleteFolder(folder string) error
	DeleteFile(key string) error
}

// S3Storage allows to interact with S3 to store files
type S3Storage struct {
	S3     *s3.S3
	Bucket string
}

// NewS3Storage creates new S3 session
func NewS3Storage(bucket, origin string) (*S3Storage, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return nil, err
	}

	client := s3.New(session)

	rule := s3.CORSRule{
		AllowedHeaders: aws.StringSlice([]string{"Authorization", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "accept", "origin", "Cache-Control", " X-Requested-With"}),
		AllowedOrigins: aws.StringSlice([]string{origin}),
		MaxAgeSeconds:  aws.Int64(3000),

		AllowedMethods: aws.StringSlice([]string{"PUT", "GET", "DELETE"}),
	}

	client.PutBucketCors(&s3.PutBucketCorsInput{
		Bucket: aws.String(bucket),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: []*s3.CORSRule{&rule},
		},
	})

	return &S3Storage{
		S3:     client,
		Bucket: bucket,
	}, nil
}

// DeleteFolder deletes every file in aws folder (prefixed: <folder>/)
func (s *S3Storage) DeleteFolder(folder string) error {

	response, err := s.S3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(folder + "/"),
	})
	if err != nil {
		return err
	}

	var objects []*s3.ObjectIdentifier
	for _, object := range response.Contents {
		objects = append(objects, &s3.ObjectIdentifier{Key: object.Key})
	}

	_, err = s.S3.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(s.Bucket),
		Delete: &s3.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// GetPresignedPutRequest creates a new PUT request and signs it with application credentials
func (s *S3Storage) GetPresignedPutRequest(key string) (string, error) {
	req, _ := s.S3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(30 * time.Second)
	if err != nil {
		return "", err
	}

	return url, nil
}

// DeleteFile deletes a file with a given key
func (s *S3Storage) DeleteFile(key string) error {
	_, err := s.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}
