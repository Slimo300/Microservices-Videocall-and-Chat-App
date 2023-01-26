package storage

import (
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	S3     *s3.S3
	Bucket string
}

func Setup(bucket string) (*S3Storage, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return nil, err
	}
	return &S3Storage{
		S3:     s3.New(session),
		Bucket: bucket,
	}, nil
}

func (s *S3Storage) UploadFile(file multipart.File, key string) error {
	_, err := s.S3.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) DeleteFile(key string) error {
	_, err := s.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

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
