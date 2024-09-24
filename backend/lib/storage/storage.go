package storage

import "context"

type Storage interface {
	GetPresignedGetRequests(ctx context.Context, files ...GetFileInput) ([]GetFileOutput, error)
	GetPresignedPutRequests(ctx context.Context, files ...PutFileInput) ([]PutFileOutput, error)
	DeleteFile(ctx context.Context, key string) error
	DeleteFilesByPrefix(ctx context.Context, prefix string) error
}

type PutFileInput struct {
	FileName string
	Prefix   string
	FileSize int64
}

type PutFileOutput struct {
	OriginalName string
	Key          string
	PresignedURL string
}

type GetFileInput string

type GetFileOutput struct {
	Key          string
	PresignedURL string
}
