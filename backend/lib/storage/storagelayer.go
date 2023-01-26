package storage

import "mime/multipart"

type StorageLayer interface {
	UploadFile(img multipart.File, key string) error
	DeleteFile(key string) error
	GetPresignedPutRequest(key string) (string, error)
}
