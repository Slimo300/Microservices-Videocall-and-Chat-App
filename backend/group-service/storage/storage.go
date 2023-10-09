package storage

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const MAX_BUCKET_SIZE = 4900000000 // 4.9GB

// StorageLayer describes Storage functionality (uploading and deleting files)
type StorageLayer interface {
	UploadFile(img multipart.File, key string) error
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

	if _, err := client.PutBucketCors(&s3.PutBucketCorsInput{
		Bucket: aws.String(bucket),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: []*s3.CORSRule{&rule},
		},
	}); err != nil {
		return nil, err
	}

	return &S3Storage{
		S3:     client,
		Bucket: bucket,
	}, nil
}

// UploadFile uploads file with a given key
func (s *S3Storage) UploadFile(file multipart.File, key string) error {

	ok, err := s.canUploadFile(file)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Can't upload file. Storage limit exceeded")
	}

	_, err = s.S3.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

// DeleteFile deletes a file with a given key
func (s *S3Storage) DeleteFile(key string) error {
	_, err := s.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) canUploadFile(file multipart.File) (bool, error) {

	fileSize, err := fileSize(file)
	if err != nil {
		return false, err
	}

	objects, err := s.S3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		return false, err
	}

	var bucketSize int64 = 0

	for _, obj := range objects.Contents {
		bucketSize += *obj.Size
	}

	if bucketSize+*fileSize > MAX_BUCKET_SIZE {
		return false, nil
	}

	return true, nil

}

func fileSize(file multipart.File) (*int64, error) {

	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return &fileSize, nil
}
