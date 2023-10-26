package storage

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const MAX_BUCKET_SIZE = 4900000000 // 4.9GB
const DEFaULT_REGION = "eu-central-1"

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

func NewS3Storage(accessKey, secretKey, bucket string) (*S3Storage, error) {

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(DEFaULT_REGION),
	}

	if os.Getenv("STORAGE_USE_DO") == "true" {
		config.Endpoint = aws.String("https://fra1.digitaloceanspaces.com")
	}
	if len(os.Getenv("STORAGE_REGION")) != 0 {
		config.Region = aws.String(os.Getenv("STORAGE_REGION"))
	}

	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		S3:     s3.New(session),
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

	log.Println(file)
	_, err = s.S3.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		ACL:    aws.String("public-read"),
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
