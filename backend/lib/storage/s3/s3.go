package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/google/uuid"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
)

// DeleteFile deletes a file with a given key
func (s *S3Storage) DeleteFile(ctx context.Context, key string) error {
	_, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// DeleteFolder deletes every file in aws folder (prefixed: <directory>/)
func (s *S3Storage) DeleteFilesByPrefix(ctx context.Context, prefix string) error {
	response, err := s.s3.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix + "/"),
	})
	if err != nil {
		return err
	}

	var objects []types.ObjectIdentifier
	for _, object := range response.Contents {
		objects = append(objects, types.ObjectIdentifier{Key: object.Key})
	}

	_, err = s.s3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Storage) GetPresignedPutRequests(ctx context.Context, files ...storage.PutFileInput) ([]storage.PutFileOutput, error) {
	ok, err := s.canUploadFiles(ctx, files...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("can't upload files")
	}

	var out []storage.PutFileOutput

	presignClient := s3.NewPresignClient(s.s3)

	for _, fileInput := range files {
		key := fileInput.Prefix + uuid.NewString()

		req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(s.bucket),
			Key:           aws.String(key),
			ContentLength: aws.Int64(fileInput.FileSize),
		}, s3.WithPresignExpires(s.presignExpiration*time.Second))
		if err != nil {
			return nil, err
		}

		out = append(out, storage.PutFileOutput{
			OriginalName: fileInput.FileName,
			Key:          key,
			PresignedURL: req.URL,
		})
	}

	return out, nil
}

func (s *S3Storage) GetPresignedGetRequests(ctx context.Context, files ...storage.GetFileInput) ([]storage.GetFileOutput, error) {
	var out []storage.GetFileOutput

	for _, key := range files {
		strKey := string(key)

		presignClient := s3.NewPresignClient(s.s3)

		req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(strKey),
		}, s3.WithPresignExpires(s.presignExpiration*time.Second))
		if err != nil {
			return nil, err
		}

		out = append(out, storage.GetFileOutput{
			Key:          strKey,
			PresignedURL: req.URL,
		})
	}

	return out, nil
}

func (s *S3Storage) canUploadFiles(ctx context.Context, filesInput ...storage.PutFileInput) (bool, error) {
	// Calculating size of files that will be uploaded
	var filesSize int64 = 0
	for _, input := range filesInput {
		if input.FileSize > s.maxFileSize {
			return false, ErrFileTooBig{fileSize: input.FileSize, fileName: input.FileName, maxFileSize: s.maxFileSize}
		}
		filesSize += input.FileSize
	}
	// Getting bucket objects list
	objects, err := s.s3.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return false, fmt.Errorf("error occured when getting objects from bucket: %w", err)
	}
	// Calculating current bucket size
	var bucketSize int64 = 0
	for _, obj := range objects.Contents {
		bucketSize += *obj.Size
	}
	if bucketSize+filesSize > s.maxBucketSize {
		return false, ErrSpaceLimitExceeded{maxUsage: s.maxBucketSize, currentUsage: bucketSize}
	}

	return true, nil
}
