package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/google/uuid"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage"
)

func (s *S3Storage) DeleteFile(ctx context.Context, key string) error {
	_, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Storage) UploadFile(ctx context.Context, key string, file multipart.File, acl storage.ACL) error {
	ok, err := s.canUploadFile(ctx, file)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("can't upload files")
	}

	_, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		ACL:    getCannedACL(acl),
	})
	return err
}

func (s *S3Storage) DeleteFilesByPrefix(ctx context.Context, prefix string) error {
	response, err := s.s3.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
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

func (s *S3Storage) GetPresignedPutRequests(ctx context.Context, inputs ...storage.PresignPutFileInput) ([]storage.PresignPutFileOutput, error) {
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
	}

	ok, err := s.canUploadFiles(ctx, inputs...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("can't upload files")
	}

	presignClient := s3.NewPresignClient(s.s3)
	var out []storage.PresignPutFileOutput
	for _, input := range inputs {
		key := input.Prefix + uuid.NewString()

		req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String(s.bucket),
			Key:           aws.String(key),
			ContentLength: aws.Int64(input.FileSize),
		}, s3.WithPresignExpires(s.presignExpiration*time.Second))
		if err != nil {
			return nil, err
		}

		out = append(out, storage.PresignPutFileOutput{
			OriginalName: input.FileName,
			Key:          key,
			PresignedURL: req.URL,
		})
	}

	return out, nil
}

func (s *S3Storage) GetPresignedGetRequests(ctx context.Context, inputs ...storage.PresignGetFileInput) ([]storage.PresignGetFileOutput, error) {
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
	}
	var out []storage.PresignGetFileOutput
	for _, input := range inputs {
		strKey := string(input)

		presignClient := s3.NewPresignClient(s.s3)
		req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(strKey),
		}, s3.WithPresignExpires(s.presignExpiration*time.Second))
		if err != nil {
			return nil, err
		}

		out = append(out, storage.PresignGetFileOutput{
			Key:          strKey,
			PresignedURL: req.URL,
		})
	}

	return out, nil
}

func (s *S3Storage) getBucketSize(ctx context.Context) (int64, error) {
	objects, err := s.s3.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return 0, fmt.Errorf("error occured when getting objects from bucket: %w", err)
	}
	var bucketSize int64 = 0
	for _, obj := range objects.Contents {
		bucketSize += *obj.Size
	}

	return bucketSize, nil
}

func (s *S3Storage) canUploadFiles(ctx context.Context, filesInput ...storage.PresignPutFileInput) (bool, error) {
	var filesSize int64 = 0
	for _, input := range filesInput {
		if input.FileSize > s.maxFileSize {
			return false, ErrFileTooBig{fileSize: input.FileSize, fileName: input.FileName, maxFileSize: s.maxFileSize}
		}
		filesSize += input.FileSize
	}
	bucketSize, err := s.getBucketSize(ctx)
	if err != nil {
		return false, err
	}
	if bucketSize+filesSize > s.maxBucketSize {
		return false, ErrSpaceLimitExceeded{maxUsage: s.maxBucketSize, currentUsage: bucketSize}
	}
	return true, nil
}

func (s *S3Storage) canUploadFile(ctx context.Context, file multipart.File) (bool, error) {
	size, err := fileSize(file)
	if err != nil {
		return false, err
	}
	if size > s.maxFileSize {
		return false, ErrFileTooBig{size, "", s.maxFileSize}
	}

	bucketSize, err := s.getBucketSize(ctx)
	if err != nil {
		return false, err
	}
	if bucketSize+size > s.maxBucketSize {
		return false, ErrSpaceLimitExceeded{maxUsage: s.maxBucketSize, currentUsage: bucketSize}
	}
	return true, nil
}

func fileSize(file multipart.File) (int64, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return 0, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	return fileSize, nil
}

func getCannedACL(acl storage.ACL) types.ObjectCannedACL {
	switch acl {
	case storage.PRIVATE:
		return types.ObjectCannedACLPrivate
	case storage.PUBLIC_READ:
		return types.ObjectCannedACLPublicRead
	default:
		panic("invalid acl")
	}
}
