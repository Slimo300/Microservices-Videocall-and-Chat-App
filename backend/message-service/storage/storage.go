package storage

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

const MAX_BUCKET_SIZE = 4900000000

// StorageLayer defines functionality expected from Storage
type StorageLayer interface {
	GetPresignedPutRequests(string, ...FileInput) ([]FileOutput, error)
	DeleteFolder(folder string) error
	DeleteFile(key string) error
}

type FileInput struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type FileOutput struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	PresignedURL string `json:"url"`
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

func (s *S3Storage) GetPresignedPutRequests(prefix string, files ...FileInput) ([]FileOutput, error) {

	ok, err := s.canUploadFiles(files...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Storage Limit Exceeded")
	}

	var out []FileOutput

	for _, fileInfo := range files {
		key := prefix + "/" + uuid.NewString()

		req, _ := s.S3.PutObjectRequest(&s3.PutObjectInput{
			Bucket:        aws.String(s.Bucket),
			Key:           aws.String(key),
			ContentLength: aws.Int64(fileInfo.Size),
		})

		url, err := req.Presign(30 * time.Second)
		if err != nil {
			return nil, err
		}

		out = append(out, FileOutput{
			Name:         fileInfo.Name,
			Key:          key,
			PresignedURL: url,
		})
	}

	return out, nil
}

// DeleteFile deletes a file with a given key
func (s *S3Storage) DeleteFile(key string) error {
	_, err := s.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) canUploadFiles(files ...FileInput) (bool, error) {

	objects, err := s.S3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		return false, err
	}

	var bucketSize int64 = 0

	// Calculating current bucket size
	for _, obj := range objects.Contents {
		bucketSize += *obj.Size
	}

	var filesSize int64 = 0

	// Calculating size of uploaded files
	for _, file := range files {
		filesSize += file.Size
	}

	if bucketSize+filesSize > MAX_BUCKET_SIZE {
		return false, nil
	}

	return true, nil

}
