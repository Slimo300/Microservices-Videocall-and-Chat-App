package storage

import "mime/multipart"

// FileStorage is an interface for file systems to fulfill
type FileStorage interface {
	UploadFile(img multipart.File, key string) error
	DeleteFile(key string) error
	GetPresignedPutRequest(key string) (string, error)
	DeleteFolder(folder string) error
}
