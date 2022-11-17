package storage

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	*s3.S3
	Bucket string
}

func Setup(bucket string) S3Storage {
	return S3Storage{
		S3: s3.New(session.Must(session.NewSession(&aws.Config{
			Region: aws.String("eu-central-1"),
		}))),
		Bucket: bucket,
	}
}

func (s *S3Storage) UpdateProfilePicture(img multipart.File, key string) error {
	_, err := s.PutObject(&s3.PutObjectInput{
		Body:   img,
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) DeleteProfilePicture(key string) error {
	_, err := s.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key + ".jpeg"),
	})
	return err
}
