package s3

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/storage"
)

// S3Storage allows to interact with S3 to store files
type S3Storage struct {
	S3     *s3.S3
	Bucket string
}

type S3Option func(*S3Storage) error

func WithCORS(origin string) S3Option {
	return func(s *S3Storage) error {
		rule := s3.CORSRule{
			AllowedHeaders: aws.StringSlice([]string{"Authorization", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Accept", "Origin", "Cache-Control", " X-Requested-With"}),
			AllowedOrigins: aws.StringSlice([]string{origin}),
			MaxAgeSeconds:  aws.Int64(3000),

			AllowedMethods: aws.StringSlice([]string{"PUT", "GET", "DELETE"}),
		}

		_, err := s.S3.PutBucketCors(&s3.PutBucketCorsInput{
			Bucket: aws.String(s.Bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: []*s3.CORSRule{&rule},
			},
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func WithACL(acl string) S3Option {
	return func(s *S3Storage) error {
		_, err := s.S3.PutBucketAcl(&s3.PutBucketAclInput{
			ACL:    aws.String(acl),
			Bucket: &s.Bucket,
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func NewS3Storage(accessKey, secretKey, bucket string, options ...S3Option) (*S3Storage, error) {

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(storage.DEFAULT_REGION),
	}

	if len(os.Getenv("STORAGE_REGION")) != 0 {
		config.Region = aws.String(os.Getenv("STORAGE_REGION"))
	}

	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	storage := &S3Storage{
		S3:     s3.New(session),
		Bucket: bucket,
	}

	for _, opt := range options {
		if err := opt(storage); err != nil {
			return nil, err
		}
	}

	return storage, nil
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

func (s *S3Storage) GetPresignedPutRequests(prefix string, files ...storage.PutFileInput) ([]storage.PutFileOutput, error) {

	ok, err := s.canUploadFiles(files...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("storage limit exceeded")
	}

	var out []storage.PutFileOutput

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

		out = append(out, storage.PutFileOutput{
			Name:         fileInfo.Name,
			Key:          key,
			PresignedURL: url,
		})
	}

	return out, nil
}

func (s *S3Storage) GetPresignedGetRequests(prefix string, files ...storage.GetFileInput) ([]storage.GetFileOutput, error) {

	var out []storage.GetFileOutput

	for _, key := range files {

		strKey := string(key)

		req, _ := s.S3.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(strKey),
		})

		url, err := req.Presign(30 * time.Second)
		if err != nil {
			return nil, err
		}

		out = append(out, storage.GetFileOutput{
			Key:          strKey,
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

func (s *S3Storage) canUploadFiles(files ...storage.PutFileInput) (bool, error) {

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

	if bucketSize+filesSize > storage.MAX_BUCKET_SIZE {
		return false, nil
	}

	return true, nil

}
