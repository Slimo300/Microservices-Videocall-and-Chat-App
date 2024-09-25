package storage

import (
	"context"
	"errors"
	"mime/multipart"
)

type Storage interface {
	UploadFile(ctx context.Context, key string, file multipart.File, acl ACL) error
	GetPresignedGetRequests(ctx context.Context, files ...PresignGetFileInput) ([]PresignGetFileOutput, error)
	GetPresignedPutRequests(ctx context.Context, files ...PresignPutFileInput) ([]PresignPutFileOutput, error)
	DeleteFile(ctx context.Context, key string) error
	DeleteFilesByPrefix(ctx context.Context, prefix string) error
}

var ErrInvalidFileName = errors.New("invalid fileName")
var ErrInvalidFileSize = errors.New("invalid fileSize")
var ErrInvalidACL = errors.New("invalid ACL")

type PresignPutFileInput struct {
	FileName string
	FileSize int64
	Prefix   string
	ACL      ACL
}

func (i PresignPutFileInput) Validate() error {
	if i.FileName == "" {
		return ErrInvalidFileName
	}
	if i.FileSize <= 0 {
		return ErrInvalidFileSize
	}
	if err := i.ACL.validate(); err != nil {
		return err
	}
	return nil
}

type PresignPutFileOutput struct {
	OriginalName string
	Key          string
	PresignedURL string
}

type PresignGetFileInput string

func (i PresignGetFileInput) Validate() error {
	if i == "" {
		return ErrInvalidFileName
	}
	return nil
}

type PresignGetFileOutput struct {
	Key          string
	PresignedURL string
}

// ACL (Access Code Lists) tells storage what permissions should be required for files to interact with them
type ACL int

func (a ACL) validate() error {
	if a != PRIVATE && a != PUBLIC_READ {
		return ErrInvalidACL
	}
	return nil
}

const (
	PRIVATE ACL = iota
	PUBLIC_READ
)
